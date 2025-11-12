'use client'

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Network, Activity, Wifi, Server, HardDrive, Search } from "lucide-react";
import { formatBytes } from "@/lib/utils";
import { useState } from "react";

export default function NetworkPage() {
  const [searchQuery, setSearchQuery] = useState("");

  const devices = [
    {
      id: "DEV-001",
      name: "Web Server 01",
      ip: "192.168.1.10",
      mac: "00:1B:44:11:3A:B7",
      type: "Server",
      status: "online",
      bandwidth: { in: 1250000000, out: 850000000 },
      uptime: "99.9%",
      lastSeen: "Active now"
    },
    {
      id: "DEV-002",
      name: "Database Server",
      ip: "192.168.1.20",
      mac: "00:1B:44:11:3A:B8",
      type: "Server",
      status: "online",
      bandwidth: { in: 950000000, out: 450000000 },
      uptime: "99.8%",
      lastSeen: "Active now"
    },
    {
      id: "DEV-003",
      name: "API Gateway",
      ip: "192.168.1.30",
      mac: "00:1B:44:11:3A:B9",
      type: "Gateway",
      status: "online",
      bandwidth: { in: 2100000000, out: 1800000000 },
      uptime: "99.9%",
      lastSeen: "Active now"
    },
    {
      id: "DEV-004",
      name: "Load Balancer",
      ip: "192.168.1.40",
      mac: "00:1B:44:11:3A:C0",
      type: "Network",
      status: "warning",
      bandwidth: { in: 1500000000, out: 1200000000 },
      uptime: "98.5%",
      lastSeen: "2 minutes ago"
    },
    {
      id: "DEV-005",
      name: "Backup Server",
      ip: "192.168.1.50",
      mac: "00:1B:44:11:3A:C1",
      type: "Server",
      status: "offline",
      bandwidth: { in: 0, out: 0 },
      uptime: "95.2%",
      lastSeen: "1 hour ago"
    },
    {
      id: "DEV-006",
      name: "Firewall",
      ip: "192.168.1.1",
      mac: "00:1B:44:11:3A:C2",
      type: "Security",
      status: "online",
      bandwidth: { in: 3200000000, out: 2900000000 },
      uptime: "99.9%",
      lastSeen: "Active now"
    }
  ];

  const getStatusColor = (status: string) => {
    switch (status) {
      case "online": return "bg-green-500";
      case "warning": return "bg-yellow-500";
      case "offline": return "bg-red-500";
      default: return "bg-gray-500";
    }
  };

  const getStatusBadge = (status: string) => {
    switch (status) {
      case "online": return "default";
      case "warning": return "default";
      case "offline": return "destructive";
      default: return "secondary";
    }
  };

  const filteredDevices = devices.filter(device =>
    device.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    device.ip.includes(searchQuery) ||
    device.type.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const onlineDevices = devices.filter(d => d.status === "online").length;
  const totalBandwidthIn = devices.reduce((sum, d) => sum + d.bandwidth.in, 0);
  const totalBandwidthOut = devices.reduce((sum, d) => sum + d.bandwidth.out, 0);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">Network Monitoring</h2>
          <p className="text-muted-foreground">
            Real-time network device monitoring and bandwidth analysis
          </p>
        </div>
        <Button>
          <Activity className="w-4 h-4 mr-2" />
          Refresh
        </Button>
      </div>

      {/* Stats */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">Total Devices</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{devices.length}</div>
            <p className="text-xs text-muted-foreground">
              {onlineDevices} online
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">Network Status</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-500">Healthy</div>
            <p className="text-xs text-muted-foreground">
              99.8% uptime
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">Bandwidth In</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatBytes(totalBandwidthIn)}/s</div>
            <p className="text-xs text-muted-foreground">
              Total incoming
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">Bandwidth Out</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatBytes(totalBandwidthOut)}/s</div>
            <p className="text-xs text-muted-foreground">
              Total outgoing
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Search */}
      <Card>
        <CardHeader>
          <CardTitle>Network Devices</CardTitle>
          <CardDescription>Monitor and manage network devices</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="relative mb-6">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
            <Input
              placeholder="Search devices..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-10"
            />
          </div>

          {/* Devices Grid */}
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {filteredDevices.map((device) => (
              <Card key={device.id} className="border-2 hover:border-primary transition-colors">
                <CardHeader className="pb-3">
                  <div className="flex items-start justify-between">
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center">
                        {device.type === "Server" && <Server className="w-5 h-5 text-primary" />}
                        {device.type === "Gateway" && <Network className="w-5 h-5 text-primary" />}
                        {device.type === "Network" && <Wifi className="w-5 h-5 text-primary" />}
                        {device.type === "Security" && <HardDrive className="w-5 h-5 text-primary" />}
                      </div>
                      <div>
                        <CardTitle className="text-base">{device.name}</CardTitle>
                        <p className="text-xs text-muted-foreground">{device.type}</p>
                      </div>
                    </div>
                    <div className={`w-2 h-2 rounded-full ${getStatusColor(device.status)}`} />
                  </div>
                </CardHeader>
                <CardContent className="space-y-3">
                  <div className="space-y-1">
                    <div className="flex justify-between text-sm">
                      <span className="text-muted-foreground">IP Address</span>
                      <span className="font-mono">{device.ip}</span>
                    </div>
                    <div className="flex justify-between text-sm">
                      <span className="text-muted-foreground">MAC Address</span>
                      <span className="font-mono text-xs">{device.mac}</span>
                    </div>
                    <div className="flex justify-between text-sm">
                      <span className="text-muted-foreground">Uptime</span>
                      <span className="font-medium">{device.uptime}</span>
                    </div>
                  </div>

                  <div className="pt-2 border-t space-y-2">
                    <div className="flex justify-between text-sm">
                      <span className="text-muted-foreground">In</span>
                      <span className="font-medium text-green-500">
                        {formatBytes(device.bandwidth.in)}/s
                      </span>
                    </div>
                    <div className="flex justify-between text-sm">
                      <span className="text-muted-foreground">Out</span>
                      <span className="font-medium text-blue-500">
                        {formatBytes(device.bandwidth.out)}/s
                      </span>
                    </div>
                  </div>

                  <div className="flex items-center justify-between pt-2">
                    <Badge variant={getStatusBadge(device.status)}>
                      {device.status}
                    </Badge>
                    <span className="text-xs text-muted-foreground">{device.lastSeen}</span>
                  </div>

                  <Button variant="outline" size="sm" className="w-full">
                    View Details
                  </Button>
                </CardContent>
              </Card>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Network Traffic Chart */}
      <Card>
        <CardHeader>
          <CardTitle>Network Traffic</CardTitle>
          <CardDescription>Real-time bandwidth usage over time</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="h-[300px] flex items-center justify-center bg-muted/50 rounded-lg">
            <p className="text-muted-foreground">Network traffic chart would go here</p>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
