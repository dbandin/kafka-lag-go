package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

func (rm *RedisManager) RegisterNode(nodeID string, ttl int) (int, error) {
	return rm.refreshNode(nodeID, ttl)
}

// DeregisterNode deregisters a node by its ID using the Lua script.
// This function ensures that the node is properly deregistered from Redis.
func (rm *RedisManager) DeregisterNode(nodeID string) error {
	_, err := rm.client.EvalSha(rm.ctx, rm.LuaSHA, []string{"deregister"}, nodeID).Result()
	if err != nil {
		return fmt.Errorf("error executing Lua script for deregistration: %v", err)
	}
	return nil
}

func (rm *RedisManager) refreshNode(nodeID string, ttl int) (int, error) {
	result, err := rm.client.EvalSha(rm.ctx, rm.LuaSHA, []string{"register_or_refresh"}, nodeID, ttl).Result()
	if err != nil {
		return -1, fmt.Errorf("error executing Lua script: %v", err)
	}

	// Interpret the result as the node's index
	nodeIndex, ok := result.(int64)
	if !ok {
		return -1, fmt.Errorf("unexpected result format: %v", result)
	}

	logrus.Debugf("Node %s TTL updated successfully with index: %d", nodeID, nodeIndex)
	return int(nodeIndex), nil
}

// GetNodeInfo retrieves information about a node, such as its index and the total number of nodes.
func (rm *RedisManager) GetNodeInfo(nodeID string) (int, int, error) {
	result, err := rm.client.EvalSha(rm.ctx, rm.LuaSHA, []string{"get_node_info"}, nodeID).Result()
	if err != nil {
		return -1, -1, fmt.Errorf("error executing Lua script for getting node info: %v", err)
	}

	// Safe type assertion
	res, ok := result.([]interface{})
	if !ok || len(res) < 1 {
		return -1, -1, fmt.Errorf("unexpected result format: %v", result)
	}

	status, ok := res[0].(string)
	if !ok || status != "ok" {
		if status == "not_found" {
			return -1, -1, fmt.Errorf("node not found: %s", nodeID)
		}
		if status == "error" && len(res) > 1 {
			return -1, -1, fmt.Errorf("error from Lua script: %v", res[1])
		}
		return -1, -1, fmt.Errorf("unexpected status in result: %v", status)
	}

	if len(res) != 3 {
		return -1, -1, fmt.Errorf("unexpected result length: %v", result)
	}

	// Extract the node index and total nodes
	index, ok := res[1].(int64)
	if !ok {
		return -1, -1, fmt.Errorf("unexpected index format in Lua script result")
	}

	totalNodes, ok := res[2].(int64)
	if !ok {
		return -1, -1, fmt.Errorf("unexpected total nodes format in Lua script result")
	}

	return int(index), int(totalNodes), nil
}

func (rm *RedisManager) StartNodeHeartbeat(nodeID string, heartbeatInterval time.Duration, ttl int) {
	go func() {
		ticker := time.NewTicker(heartbeatInterval)
		defer ticker.Stop()

		for {
			select {
			case <-rm.ctx.Done(): // Graceful shutdown if the context is canceled
				logrus.Warnf("Node heartbeat stopped: %v", rm.ctx.Err())
				return
			case <-ticker.C: // Triggered at every tick interval
				_, err := rm.refreshNode(nodeID, ttl)
				if err != nil {
					logrus.Warnf("Failed to refresh node registration: %v", err)
				}
			}
		}
	}()
}

func (rm *RedisManager) StartNodeMonitoring(monitorInterval time.Duration) {
	go func() {
		ticker := time.NewTicker(monitorInterval)
		defer ticker.Stop()
		for {
			select {
			case <-rm.ctx.Done():
				logrus.Warnf("Node monitoring stopped: %v", rm.ctx.Err())
				return
			case <-ticker.C:
				failedNodes, err := rm.monitorNodes()
				if err != nil {
					logrus.Warnf("Failed to monitor nodes: %v", err)
				}
				if len(failedNodes) > 0 {
					logrus.Debugf("Removed failed nodes: %v", failedNodes)
				}
			}
		}
	}()
}

func (rm *RedisManager) monitorNodes() ([]string, error) {
	result, err := rm.client.EvalSha(rm.ctx, rm.LuaSHA, []string{"monitor"}).Result()
	if err != nil {
		return nil, fmt.Errorf("error executing Lua script for monitor nodes: %v", err)
	}

	failedNodes, ok := result.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	failedNodeIDs := make([]string, len(failedNodes))
	for i, nodeID := range failedNodes {
		strID, ok := nodeID.(string)
		if !ok {
			return nil, fmt.Errorf("unexpected node ID type: %T", nodeID)
		}
		failedNodeIDs[i] = strID
	}

	return failedNodeIDs, nil
}

// Ping checks the connection to Redis by sending a PING command.
func (rm *RedisManager) Ping(ctx ...context.Context) error {
	// Determine which context to use: the provided one or the default rm.ctx
	var useCtx context.Context
	if len(ctx) > 0 && ctx[0] != nil {
		useCtx = ctx[0]
	} else {
		useCtx = rm.ctx
	}

	// Use the selected context for the Ping operation
	_, err := rm.client.Ping(useCtx).Result()
	if err != nil {
		logrus.Errorf("Redis Ping failed: %v", err)
		return err
	}
	return nil
}
