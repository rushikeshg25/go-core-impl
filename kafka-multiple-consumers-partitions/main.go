package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// ============================================================================
// WHAT IS ZOOKEEPER? 🤔
// ============================================================================

/*
🏛️ ZOOKEEPER EXPLAINED:

Think of Zookeeper like the "Manager" of your Kafka cluster:

1. 📋 CLUSTER COORDINATION:
   - Keeps track of which brokers are alive
   - Manages cluster membership
   - Handles leader election for partitions

2. 📚 METADATA STORAGE:
   - Stores topic configurations
   - Tracks partition locations
   - Manages consumer group information

3. 🗳️ LEADER ELECTION:
   - Decides which broker leads each partition
   - Handles failover when brokers go down
   - Ensures data consistency

4. 🔧 CONFIGURATION MANAGEMENT:
   - Stores dynamic configurations
   - Manages ACLs (Access Control Lists)
   - Handles quota configurations

📈 EVOLUTION:
- OLD WAY: Kafka + Zookeeper (what most tutorials show)
- NEW WAY: KRaft mode (Kafka without Zookeeper) - Kafka 2.8+
- FUTURE: Zookeeper will be completely removed

🎯 BOTH CONFLUENT AND APACHE IMAGES SUPPORT BOTH MODES!
*/

// ============================================================================
// ADVANCED KAFKA CONCEPTS - NEXT LEVEL 🚀
// ============================================================================

// Advanced message structure with headers and metadata
type AdvancedMessage struct {
	ID            string                 `json:"id"`
	EventType     string                 `json:"event_type"`
	UserID        string                 `json:"user_id"`
	Payload       map[string]interface{} `json:"payload"`
	Timestamp     time.Time              `json:"timestamp"`
	Version       string                 `json:"version"`
	Source        string                 `json:"source"`
	CorrelationID string                 `json:"correlation_id"`
}

// Message headers for tracing and metadata
type MessageHeaders struct {
	CorrelationID string
	Source        string
	ContentType   string
	Version       string
	TraceID       string
}

// ============================================================================
// CONCEPT 1: CONSUMER GROUPS & PARTITIONS 👥
// ============================================================================

func demonstrateConsumerGroups() {
	fmt.Println("🔍 CONSUMER GROUPS EXPLAINED:")
	fmt.Println()
	fmt.Println("Imagine a pizza delivery system:")
	fmt.Println("📦 Topic = 'pizza-orders'")
	fmt.Println("🍕 Partitions = Different neighborhoods (North, South, East, West)")
	fmt.Println("🚗 Consumers = Delivery drivers")
	fmt.Println("👥 Consumer Group = All drivers working for the same company")
	fmt.Println()
	fmt.Println("RULES:")
	fmt.Println("✅ Each partition can only be read by ONE consumer in a group")
	fmt.Println("✅ Multiple groups can read the SAME partition independently")
	fmt.Println("✅ If you add more consumers than partitions, some will be idle")
	fmt.Println()

	// This is why partitions matter for scaling!
	demoPartitionedProducer()
}

func demoPartitionedProducer() {
	config := kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	}

	producer, err := kafka.NewProducer(&config)
	if err != nil {
		log.Printf("Failed to create producer: %v", err)
		return
	}
	defer producer.Close()

	topic := "advanced-topic"

	fmt.Println("📤 SENDING MESSAGES TO DIFFERENT PARTITIONS:")

	// Send messages to specific partitions
	for i := 0; i < 6; i++ {
		partition := int32(i % 3) // Round-robin across 3 partitions

		message := AdvancedMessage{
			ID:            fmt.Sprintf("msg-%d", i),
			EventType:     "user_action",
			UserID:        fmt.Sprintf("user-%d", i%2), // Two users
			Payload:       map[string]interface{}{"action": "click", "page": "home"},
			Timestamp:     time.Now(),
			Version:       "1.0",
			Source:        "web-app",
			CorrelationID: fmt.Sprintf("corr-%d", i),
		}

		jsonBytes, _ := json.Marshal(message)

		kafkaMsg := &kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: partition, // Specific partition
			},
			Key:   []byte(message.UserID), // Same key = same partition
			Value: jsonBytes,
			Headers: []kafka.Header{
				{Key: "correlation-id", Value: []byte(message.CorrelationID)},
				{Key: "content-type", Value: []byte("application/json")},
				{Key: "version", Value: []byte(message.Version)},
			},
		}

		err := producer.Produce(kafkaMsg, nil)
		if err != nil {
			log.Printf("Failed to produce message: %v", err)
		} else {
			fmt.Printf("✅ Sent message %d to partition %d (user: %s)\n", i, partition, message.UserID)
		}
	}

	// Wait for message deliveries
	producer.Flush(15 * 1000)
}

// ============================================================================
// CONCEPT 2: CONSUMER GROUPS IN ACTION 🎭
// ============================================================================

func runAdvancedConsumerGroup(groupID string, consumerID string) {
	config := kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          groupID,
		"auto.offset.reset": "earliest",

		// Advanced consumer settings
		"enable.auto.commit":    false, // Manual commit for reliability
		"session.timeout.ms":    6000,
		"heartbeat.interval.ms": 2000,
		"max.poll.interval.ms":  300000,
		"fetch.min.bytes":       1024,
		"fetch.wait.max.ms":     500,
	}

	consumer, err := kafka.NewConsumer(&config)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	topics := []string{"advanced-topic"}
	err = consumer.SubscribeTopics(topics, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	fmt.Printf("🚀 Consumer %s (Group: %s) started\n", consumerID, groupID)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case sig := <-sigchan:
			fmt.Printf("🛑 Consumer %s caught signal %v: terminating\n", consumerID, sig)
			return

		default:
			msg, err := consumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				// Timeout is normal, continue
				continue
			}

			// Process the message
			var advMsg AdvancedMessage
			if err := json.Unmarshal(msg.Value, &advMsg); err != nil {
				log.Printf("❌ Failed to unmarshal: %v", err)
				continue
			}

			fmt.Printf("📨 [%s] Received from partition %d, offset %d:\n",
				consumerID, msg.TopicPartition.Partition, msg.TopicPartition.Offset)
			fmt.Printf("   🎯 Event: %s, User: %s, CorrelationID: %s\n",
				advMsg.EventType, advMsg.UserID, advMsg.CorrelationID)

			// Process headers
			for _, header := range msg.Headers {
				fmt.Printf("   📋 Header %s: %s\n", header.Key, string(header.Value))
			}

			// IMPORTANT: Manual commit for reliability
			_, err = consumer.Commit()
			if err != nil {
				log.Printf("❌ Failed to commit: %v", err)
			}
		}
	}
}

// ============================================================================
// CONCEPT 3: ADVANCED PATTERNS 🏗️
// ============================================================================

// Pattern 1: Exactly-Once Processing
func exactlyOnceProducer() {
	config := kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",

		// Exactly-once semantics
		"enable.idempotence":     true,
		"transactional.id":       "my-transactional-producer",
		"transaction.timeout.ms": 30000,

		// Strong durability guarantees
		"acks":                                  "all",
		"retries":                               2147483647,
		"max.in.flight.requests.per.connection": 5,
	}

	producer, err := kafka.NewProducer(&config)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	// Initialize transactions
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = producer.InitTransactions(ctx)
	if err != nil {
		log.Fatalf("Failed to init transactions: %v", err)
	}

	fmt.Println("💎 EXACTLY-ONCE PRODUCER STARTED")

	// Begin transaction
	err = producer.BeginTransaction()
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}

	topic := "exactly-once-topic"

	// Send multiple messages in one transaction
	for i := 0; i < 3; i++ {
		message := &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Key:            []byte(fmt.Sprintf("key-%d", i)),
			Value:          []byte(fmt.Sprintf("Exactly-once message %d", i)),
		}

		err = producer.Produce(message, nil)
		if err != nil {
			producer.AbortTransaction(ctx)
			log.Fatalf("Failed to produce: %v", err)
		}
	}

	// Commit transaction - all messages delivered atomically
	err = producer.CommitTransaction(ctx)
	if err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	fmt.Println("✅ Transaction committed - all 3 messages delivered exactly once!")
}

// Pattern 2: Dead Letter Queue
func deadLetterQueuePattern() {
	fmt.Println("⚰️ DEAD LETTER QUEUE PATTERN:")
	fmt.Println("When message processing fails repeatedly:")
	fmt.Println("1️⃣ Try to process message")
	fmt.Println("2️⃣ If it fails, retry up to N times")
	fmt.Println("3️⃣ After N failures, send to Dead Letter Queue")
	fmt.Println("4️⃣ Alert operations team")
	fmt.Println("5️⃣ Investigate and reprocess manually")
}

// Pattern 3: Saga Pattern (Distributed Transactions)
func sagaPattern() {
	fmt.Println("🎭 SAGA PATTERN (Distributed Transactions):")
	fmt.Println("For operations across multiple services:")
	fmt.Println("1️⃣ Order Service: Create order → Success")
	fmt.Println("2️⃣ Payment Service: Charge card → Success")
	fmt.Println("3️⃣ Inventory Service: Reserve items → FAILURE!")
	fmt.Println("4️⃣ Compensation: Refund payment, cancel order")
	fmt.Println("💡 Each step publishes events, failures trigger compensations")
}

// ============================================================================
// CONCEPT 4: MONITORING & OBSERVABILITY 📊
// ============================================================================

func monitoringAndMetrics() {
	fmt.Println("📊 PRODUCTION MONITORING:")
	fmt.Println()
	fmt.Println("KEY METRICS TO TRACK:")
	fmt.Println("🔄 Producer:")
	fmt.Println("  - Messages per second")
	fmt.Println("  - Error rate")
	fmt.Println("  - Latency (produce time)")
	fmt.Println("  - Batch size efficiency")
	fmt.Println()
	fmt.Println("📥 Consumer:")
	fmt.Println("  - Consumer lag (how far behind)")
	fmt.Println("  - Processing time per message")
	fmt.Println("  - Error rate")
	fmt.Println("  - Rebalance frequency")
	fmt.Println()
	fmt.Println("🏛️ Broker:")
	fmt.Println("  - Disk usage")
	fmt.Println("  - Network throughput")
	fmt.Println("  - Partition count")
	fmt.Println("  - Under-replicated partitions")
}

// ============================================================================
// CONCEPT 5: SCHEMA EVOLUTION 🧬
// ============================================================================

// Version 1 of user event
type UserEventV1 struct {
	UserID    string `json:"user_id"`
	Action    string `json:"action"`
	Timestamp int64  `json:"timestamp"`
}

// Version 2 - Added optional fields (backward compatible)
type UserEventV2 struct {
	UserID        string            `json:"user_id"`
	Action        string            `json:"action"`
	Timestamp     int64             `json:"timestamp"`
	SessionID     string            `json:"session_id,omitempty"` // Optional
	Properties    map[string]string `json:"properties,omitempty"` // Optional
	SchemaVersion int               `json:"schema_version"`       // Track version
}

func schemaEvolutionDemo() {
	fmt.Println("🧬 SCHEMA EVOLUTION:")
	fmt.Println("How to change message structure without breaking existing consumers:")
	fmt.Println()
	fmt.Println("✅ SAFE CHANGES (Backward Compatible):")
	fmt.Println("  - Add optional fields")
	fmt.Println("  - Add default values")
	fmt.Println("  - Make required fields optional")
	fmt.Println()
	fmt.Println("❌ UNSAFE CHANGES (Breaking):")
	fmt.Println("  - Remove fields")
	fmt.Println("  - Change field types")
	fmt.Println("  - Make optional fields required")
	fmt.Println()
	fmt.Println("🎯 SOLUTION: Schema Registry (Confluent feature)")
	fmt.Println("  - Enforces schema compatibility")
	fmt.Println("  - Version management")
	fmt.Println("  - Avro/JSON Schema support")
}

// ============================================================================
// MAIN FUNCTION - CHOOSE YOUR ADVENTURE! 🗺️
// ============================================================================

func main() {
	if len(os.Args) < 2 {
		fmt.Println("🚀 ADVANCED KAFKA CONCEPTS")
		fmt.Println("=========================")
		fmt.Println()
		fmt.Println("You've mastered the basics! Now let's go advanced:")
		fmt.Println()
		fmt.Println("CONCEPTS TO EXPLORE:")
		fmt.Println("🔍 consumer-groups     - Learn about scaling consumers")
		fmt.Println("📤 partitioned-producer - Send to specific partitions")
		fmt.Println("👥 multi-consumer      - Run multiple consumers")
		fmt.Println("💎 exactly-once        - Transactional messaging")
		fmt.Println("⚰️  patterns           - Advanced patterns (DLQ, Saga)")
		fmt.Println("📊 monitoring          - Production monitoring")
		fmt.Println("🧬 schema              - Schema evolution")
		fmt.Println()
		fmt.Println("Usage: go run main.go <concept>")
		fmt.Println()
		fmt.Println("🎯 START WITH: go run main.go consumer-groups")
		os.Exit(1)
	}

	concept := os.Args[1]

	switch concept {
	case "consumer-groups":
		demonstrateConsumerGroups()

	case "partitioned-producer":
		demoPartitionedProducer()

	case "multi-consumer":
		if len(os.Args) < 4 {
			fmt.Println("Usage: go run main.go multi-consumer <group-id> <consumer-id>")
			fmt.Println("Example: go run main.go multi-consumer pizza-delivery driver-1")
			os.Exit(1)
		}
		groupID := os.Args[2]
		consumerID := os.Args[3]
		runAdvancedConsumerGroup(groupID, consumerID)

	case "exactly-once":
		exactlyOnceProducer()

	case "patterns":
		deadLetterQueuePattern()
		fmt.Println()
		sagaPattern()

	case "monitoring":
		monitoringAndMetrics()

	case "schema":
		schemaEvolutionDemo()

	default:
		fmt.Printf("❌ Unknown concept: %s\n", concept)
		fmt.Println("Run without arguments to see available concepts")
	}
}

/*
🎯 NEXT LEVEL CONCEPTS SUMMARY:

1️⃣ CONSUMER GROUPS:
   - Multiple consumers work together
   - Each partition assigned to one consumer
   - Automatic rebalancing when consumers join/leave

2️⃣ PARTITIONING STRATEGY:
   - Same key → same partition (ordering guaranteed)
   - Round-robin for load balancing
   - Custom partitioners for special needs

3️⃣ ADVANCED PATTERNS:
   - Exactly-once processing (transactions)
   - Dead letter queues (error handling)
   - Saga pattern (distributed transactions)
   - Event sourcing (audit trail)

4️⃣ PRODUCTION CONCERNS:
   - Consumer lag monitoring
   - Schema evolution
   - Security (SASL/SSL)
   - Multi-datacenter replication

5️⃣ PERFORMANCE TUNING:
   - Batch size optimization
   - Compression settings
   - Memory and network tuning
   - Partition count planning

🚀 WORKS WITH BOTH:
- ✅ Apache Kafka Docker image
- ✅ Confluent Platform Docker image
- ✅ Any Kafka cluster (cloud or on-premise)

The concepts are universal - only the setup differs!
*/
