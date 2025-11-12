'use client'

import { useEffect, useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { AttackFlowVisualization } from "@/components/attack-flow-visualization";
import { 
  AlertTriangle, 
  Shield, 
  Activity, 
  TrendingUp,
  ArrowUpRight,
  ArrowDownRight,
  Clock
} from "lucide-react";
import { rustApi, WS_URLS } from "@/lib/api";
import { useSimpleWebSocket } from "@/hooks/useWebSocket";

export default function DashboardPage() {
  const [realtimeStats, setRealtimeStats] = useState({
    packets_per_second: 0,
    bytes_per_second: 0,
    active_connections: 0,
    threats_detected: 0,
  });

  // WebSocket for real-time updates
  const { data: wsStats } = useSimpleWebSocket(WS_URLS.stats);
  const { data: wsThreat } = useSimpleWebSocket(WS_URLS.threats);

  // Fetch initial stats
  useEffect(() => {
    const fetchStats = async () => {
      try {
        const response = await rustApi.getRealtimeStats();
        setRealtimeStats(response.data);
      } catch (error) {
        console.error('Failed to fetch stats:', error);
      }
    };

    fetchStats();
    const interval = setInterval(fetchStats, 5000);
    return () => clearInterval(interval);
  }, []);

  // Update stats from WebSocket
  useEffect(() => {
    if (wsStats) {
      setRealtimeStats(wsStats);
    }
  }, [wsStats]);

  const stats = [
    {
      title: "Active Connections",
      value: (realtimeStats?.active_connections || 0).toString(),
      change: "+12%",
      trend: "up",
      icon: Activity,
      color: "text-blue-500"
    },
    {
      title: "Threats Detected",
      value: (realtimeStats?.threats_detected || 0).toString(),
      change: "+8%",
      trend: "up",
      icon: Shield,
      color: "text-red-500"
    },
    {
      title: "Packets/sec",
      value: (realtimeStats?.packets_per_second || 0).toFixed(1),
      change: "+0.1%",
      trend: "up",
      icon: TrendingUp,
      color: "text-green-500"
    },
    {
      title: "Bandwidth",
      value: `${((realtimeStats?.bytes_per_second || 0) / 1024 / 1024).toFixed(2)} MB/s`,
      change: "Stable",
      trend: "neutral",
      icon: AlertTriangle,
      color: "text-purple-500"
    }
  ];

  const recentAlerts = [
    {
      id: 1,
      title: "Suspicious Login Attempt",
      severity: "high",
      time: "2 minutes ago",
      status: "active"
    },
    {
      id: 2,
      title: "Unusual Network Traffic",
      severity: "medium",
      time: "15 minutes ago",
      status: "investigating"
    },
    {
      id: 3,
      title: "Failed Authentication",
      severity: "low",
      time: "1 hour ago",
      status: "resolved"
    },
    {
      id: 4,
      title: "Port Scan Detected",
      severity: "high",
      time: "2 hours ago",
      status: "active"
    }
  ];

  const topThreats = [
    { name: "Malware.Generic", count: 145, severity: "critical" },
    { name: "Phishing.Email", count: 89, severity: "high" },
    { name: "Brute.Force", count: 67, severity: "medium" },
    { name: "SQL.Injection", count: 34, severity: "high" },
    { name: "XSS.Attack", count: 23, severity: "medium" }
  ];

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case "critical": return "bg-red-500";
      case "high": return "bg-orange-500";
      case "medium": return "bg-yellow-500";
      case "low": return "bg-blue-500";
      default: return "bg-gray-500";
    }
  };

  return (
    <div className="space-y-6">
      {/* Stats Grid */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {stats.map((stat, index) => (
          <Card key={index}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">
                {stat.title}
              </CardTitle>
              <stat.icon className={`w-4 h-4 ${stat.color}`} />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stat.value}</div>
              <p className="text-xs text-muted-foreground flex items-center gap-1 mt-1">
                {stat.trend === "up" ? (
                  <ArrowUpRight className="w-3 h-3 text-green-500" />
                ) : stat.trend === "down" ? (
                  <ArrowDownRight className="w-3 h-3 text-red-500" />
                ) : null}
                {stat.change} from last month
              </p>
            </CardContent>
          </Card>
        ))}
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
        {/* Recent Alerts */}
        <Card className="col-span-4">
          <CardHeader>
            <CardTitle>Recent Alerts</CardTitle>
            <CardDescription>
              Latest security alerts and incidents
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {recentAlerts.map((alert) => (
                <div
                  key={alert.id}
                  className="flex items-center justify-between p-4 border rounded-lg hover:bg-accent transition-colors"
                >
                  <div className="flex items-center gap-4">
                    <div className={`w-2 h-2 rounded-full ${getSeverityColor(alert.severity)}`} />
                    <div>
                      <p className="font-medium">{alert.title}</p>
                      <p className="text-sm text-muted-foreground flex items-center gap-1">
                        <Clock className="w-3 h-3" />
                        {alert.time}
                      </p>
                    </div>
                  </div>
                  <Badge variant={alert.status === "active" ? "destructive" : "secondary"}>
                    {alert.status}
                  </Badge>
                </div>
              ))}
            </div>
            <Button variant="outline" className="w-full mt-4">
              View All Alerts
            </Button>
          </CardContent>
        </Card>

        {/* Top Threats */}
        <Card className="col-span-3">
          <CardHeader>
            <CardTitle>Top Threats</CardTitle>
            <CardDescription>
              Most detected threats this week
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {topThreats.map((threat, index) => (
                <div key={index} className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <div className="text-sm font-medium text-muted-foreground">
                      #{index + 1}
                    </div>
                    <div>
                      <p className="text-sm font-medium">{threat.name}</p>
                      <p className="text-xs text-muted-foreground">
                        {threat.count} detections
                      </p>
                    </div>
                  </div>
                  <Badge variant="outline" className={getSeverityColor(threat.severity)}>
                    {threat.severity}
                  </Badge>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Attack Flow Visualization */}
      <AttackFlowVisualization />
    </div>
  );
}
