library(ggplot2)

# data for the benchmark
categories <- c("Raw eBPF", "Cilium")
throughput <- c(9500, 8500)  # Packets per second
latency <- c(10, 15)         # Latency in milliseconds
cpu_usage <- c(70, 85)       # CPU usage in percentage

# create data frames for each metric
throughput_data <- data.frame(Category = categories, Throughput = throughput)
latency_data <- data.frame(Category = categories, Latency = latency)
cpu_usage_data <- data.frame(Category = categories, CPU_Usage = cpu_usage)

# Throughput graph
throughput_plot <- ggplot(throughput_data, aes(x = Category, y = Throughput, fill = Category)) +
  geom_bar(stat = "identity") +
  labs(title = "Throughput Comparison", y = "Packets/s", x = "") +
  theme_minimal() +
  theme(legend.position = "none")

# Latency graph
latency_plot <- ggplot(latency_data, aes(x = Category, y = Latency, fill = Category)) +
  geom_bar(stat = "identity") +
  labs(title = "Latency Comparison", y = "Latency (ms)", x = "") +
  theme_minimal() +
  theme(legend.position = "none")

# CPU Usage graph
cpu_usage_plot <- ggplot(cpu_usage_data, aes(x = Category, y = CPU_Usage, fill = Category)) +
  geom_bar(stat = "identity") +
  labs(title = "CPU Usage Comparison", y = "CPU Usage (%)", x = "") +
  theme_minimal() +
  theme(legend.position = "none")

# plot all graphs together using gridExtra
library(gridExtra)
grid.arrange(throughput_plot, latency_plot, cpu_usage_plot, nrow = 1)
