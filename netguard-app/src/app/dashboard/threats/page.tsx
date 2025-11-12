'use client'

import { useEffect, useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Flame, Shield, Ban, Eye, Loader2 } from "lucide-react";
import { formatDate } from "@/lib/utils";
import { gateway } from "@/lib/api";
import { useSimpleWebSocket } from "@/hooks/useWebSocket";
import { WS_URLS } from "@/lib/api";

export default function ThreatsPage() {
  const [threats, setThreats] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  
  // WebSocket for real-time threat updates
  const { data: newThreat } = useSimpleWebSocket(WS_URLS.threats);

  // Fetch threats from API
  useEffect(() => {
    const fetchThreats = async () => {
      try {
        const response = await gateway.listThreats();
        setThreats(response.data.threats || []);
      } catch (error) {
        console.error('Failed to fetch threats:', error);
        // Use fallback data if API fails
        setThreats(fallbackThreats);
      } finally {
        setLoading(false);
      }
    };

    fetchThreats();
  }, []);

  // Add new threats from WebSocket
  useEffect(() => {
    if (newThreat) {
      setThreats((prev) => [
        {
          id: `THR-${Date.now()}`,
          name: newThreat.threat_type,
          type: newThreat.threat_type,
          severity: newThreat.severity,
          status: 'monitoring',
          source: newThreat.source_ip,
          target: newThreat.destination_ip,
          timestamp: newThreat.timestamp,
          detections: 1
        },
        ...prev
      ].slice(0, 50));
    }
  }, [newThreat]);

  const fallbackThreats = [
    {
      id: "THR-001",
      name: "Malware.Generic.Trojan",
      type: "Malware",
      severity: "critical",
      status: "blocked",
      source: "192.168.1.100",
      target: "Server-01",
      timestamp: new Date().toISOString(),
      detections: 145
    },
    {
      id: "THR-002",
      name: "Phishing.Email.Suspicious",
      type: "Phishing",
      severity: "high",
      status: "quarantined",
      source: "external@malicious.com",
      target: "Email Gateway",
      timestamp: new Date(Date.now() - 1800000).toISOString(),
      detections: 89
    },
    {
      id: "THR-003",
      name: "Brute.Force.SSH",
      type: "Brute Force",
      severity: "medium",
      status: "monitoring",
      source: "203.0.113.45",
      target: "SSH Service",
      timestamp: new Date(Date.now() - 3600000).toISOString(),
      detections: 67
    },
    {
      id: "THR-004",
      name: "SQL.Injection.Attempt",
      type: "Injection",
      severity: "high",
      status: "blocked",
      source: "198.51.100.23",
      target: "Web Application",
      timestamp: new Date(Date.now() - 7200000).toISOString(),
      detections: 34
    },
    {
      id: "THR-005",
      name: "XSS.Attack.Reflected",
      type: "XSS",
      severity: "medium",
      status: "blocked",
      source: "203.0.113.89",
      target: "Web Server",
      timestamp: new Date(Date.now() - 10800000).toISOString(),
      detections: 23
    }
  ];

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case "critical": return "destructive";
      case "high": return "destructive";
      case "medium": return "default";
      case "low": return "secondary";
      default: return "secondary";
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case "blocked": return <Ban className="w-4 h-4 text-red-500" />;
      case "quarantined": return <Shield className="w-4 h-4 text-orange-500" />;
      case "monitoring": return <Eye className="w-4 h-4 text-blue-500" />;
      default: return <Flame className="w-4 h-4 text-yellow-500" />;
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">Threat Detection</h2>
          <p className="text-muted-foreground">
            AI-powered threat detection and analysis
          </p>
        </div>
        <Button>
          <Flame className="w-4 h-4 mr-2" />
          Run Scan
        </Button>
      </div>

      {/* Stats */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">Total Threats</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{threats.length}</div>
            <p className="text-xs text-muted-foreground">Last 24 hours</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">Blocked</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-500">
              {threats.filter(t => t.status === "blocked").length}
            </div>
            <p className="text-xs text-muted-foreground">Automatically blocked</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">Quarantined</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-orange-500">
              {threats.filter(t => t.status === "quarantined").length}
            </div>
            <p className="text-xs text-muted-foreground">Under review</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">Detection Rate</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-500">99.8%</div>
            <p className="text-xs text-muted-foreground">Success rate</p>
          </CardContent>
        </Card>
      </div>

      {/* Threats Table */}
      <Card>
        <CardHeader>
          <CardTitle>Detected Threats</CardTitle>
          <CardDescription>Real-time threat detection and response</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="border rounded-lg">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>ID</TableHead>
                  <TableHead>Threat Name</TableHead>
                  <TableHead>Type</TableHead>
                  <TableHead>Severity</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Source</TableHead>
                  <TableHead>Target</TableHead>
                  <TableHead>Detections</TableHead>
                  <TableHead>Time</TableHead>
                  <TableHead>Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {threats.map((threat) => (
                  <TableRow key={threat.id}>
                    <TableCell className="font-medium">{threat.id}</TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        <Flame className="w-4 h-4 text-orange-500" />
                        <span className="font-medium">{threat.name}</span>
                      </div>
                    </TableCell>
                    <TableCell>{threat.type}</TableCell>
                    <TableCell>
                      <Badge variant={getSeverityColor(threat.severity)}>
                        {threat.severity}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        {getStatusIcon(threat.status)}
                        <span className="capitalize">{threat.status}</span>
                      </div>
                    </TableCell>
                    <TableCell className="font-mono text-sm">{threat.source}</TableCell>
                    <TableCell>{threat.target}</TableCell>
                    <TableCell>
                      <Badge variant="outline">{threat.detections}</Badge>
                    </TableCell>
                    <TableCell className="text-sm text-muted-foreground">
                      {formatDate(threat.timestamp)}
                    </TableCell>
                    <TableCell>
                      <div className="flex gap-2">
                        <Button variant="ghost" size="sm">
                          Details
                        </Button>
                        <Button variant="ghost" size="sm">
                          Block
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>

      {/* Threat Distribution */}
      <div className="grid gap-4 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Threat Types</CardTitle>
            <CardDescription>Distribution by threat category</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {["Malware", "Phishing", "Brute Force", "Injection", "XSS"].map((type, index) => (
                <div key={index} className="flex items-center justify-between">
                  <span className="text-sm font-medium">{type}</span>
                  <div className="flex items-center gap-2">
                    <div className="w-32 h-2 bg-muted rounded-full overflow-hidden">
                      <div 
                        className="h-full bg-primary" 
                        style={{ width: `${Math.random() * 100}%` }}
                      />
                    </div>
                    <span className="text-sm text-muted-foreground w-12 text-right">
                      {Math.floor(Math.random() * 100)}%
                    </span>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Top Attack Sources</CardTitle>
            <CardDescription>Most active threat sources</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {[
                { ip: "203.0.113.45", country: "Unknown", attacks: 234 },
                { ip: "198.51.100.23", country: "Unknown", attacks: 189 },
                { ip: "192.0.2.100", country: "Unknown", attacks: 156 },
                { ip: "203.0.113.89", country: "Unknown", attacks: 123 },
                { ip: "198.51.100.67", country: "Unknown", attacks: 98 }
              ].map((source, index) => (
                <div key={index} className="flex items-center justify-between p-3 border rounded-lg">
                  <div>
                    <p className="font-mono text-sm font-medium">{source.ip}</p>
                    <p className="text-xs text-muted-foreground">{source.country}</p>
                  </div>
                  <Badge variant="destructive">{source.attacks} attacks</Badge>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
