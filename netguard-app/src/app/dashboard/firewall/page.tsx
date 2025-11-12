'use client'

import { useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Shield, Plus, Search, Trash2, Edit, Power, PowerOff } from "lucide-react";

export default function FirewallPage() {
  const [searchQuery, setSearchQuery] = useState("");
  const [actionFilter, setActionFilter] = useState("all");

  const rules = [
    {
      id: "FW-001",
      name: "Allow HTTP Traffic",
      source: "0.0.0.0/0",
      destination: "192.168.1.0/24",
      port: "80",
      protocol: "TCP",
      action: "allow",
      status: "enabled",
      priority: 100,
      hits: 15234
    },
    {
      id: "FW-002",
      name: "Allow HTTPS Traffic",
      source: "0.0.0.0/0",
      destination: "192.168.1.0/24",
      port: "443",
      protocol: "TCP",
      action: "allow",
      status: "enabled",
      priority: 100,
      hits: 28456
    },
    {
      id: "FW-003",
      name: "Block Suspicious IPs",
      source: "203.0.113.0/24",
      destination: "Any",
      port: "Any",
      protocol: "Any",
      action: "deny",
      status: "enabled",
      priority: 50,
      hits: 892
    },
    {
      id: "FW-004",
      name: "Allow SSH Admin",
      source: "10.0.0.0/8",
      destination: "192.168.1.10",
      port: "22",
      protocol: "TCP",
      action: "allow",
      status: "enabled",
      priority: 90,
      hits: 456
    },
    {
      id: "FW-005",
      name: "Block Outbound Malware",
      source: "192.168.1.0/24",
      destination: "198.51.100.0/24",
      port: "Any",
      protocol: "Any",
      action: "deny",
      status: "enabled",
      priority: 60,
      hits: 234
    },
    {
      id: "FW-006",
      name: "Allow DNS",
      source: "192.168.1.0/24",
      destination: "8.8.8.8",
      port: "53",
      protocol: "UDP",
      action: "allow",
      status: "enabled",
      priority: 100,
      hits: 45678
    },
    {
      id: "FW-007",
      name: "Test Rule",
      source: "Any",
      destination: "Any",
      port: "8080",
      protocol: "TCP",
      action: "allow",
      status: "disabled",
      priority: 10,
      hits: 0
    }
  ];

  const getActionColor = (action: string) => {
    return action === "allow" ? "default" : "destructive";
  };

  const getStatusIcon = (status: string) => {
    return status === "enabled" ? (
      <Power className="w-4 h-4 text-green-500" />
    ) : (
      <PowerOff className="w-4 h-4 text-gray-400" />
    );
  };

  const filteredRules = rules.filter(rule => {
    const matchesSearch = rule.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
                         rule.source.includes(searchQuery) ||
                         rule.destination.includes(searchQuery);
    const matchesAction = actionFilter === "all" || rule.action === actionFilter;
    return matchesSearch && matchesAction;
  });

  const enabledRules = rules.filter(r => r.status === "enabled").length;
  const allowRules = rules.filter(r => r.action === "allow").length;
  const denyRules = rules.filter(r => r.action === "deny").length;
  const totalHits = rules.reduce((sum, r) => sum + r.hits, 0);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">Firewall Management</h2>
          <p className="text-muted-foreground">
            Configure and manage firewall rules
          </p>
        </div>
        <Button>
          <Plus className="w-4 h-4 mr-2" />
          Add Rule
        </Button>
      </div>

      {/* Stats */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">Total Rules</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{rules.length}</div>
            <p className="text-xs text-muted-foreground">
              {enabledRules} enabled
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">Allow Rules</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-500">{allowRules}</div>
            <p className="text-xs text-muted-foreground">
              Permissive rules
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">Deny Rules</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-500">{denyRules}</div>
            <p className="text-xs text-muted-foreground">
              Restrictive rules
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">Total Hits</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{totalHits.toLocaleString()}</div>
            <p className="text-xs text-muted-foreground">
              Rule matches
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Firewall Rules */}
      <Card>
        <CardHeader>
          <CardTitle>Firewall Rules</CardTitle>
          <CardDescription>Manage network access control rules</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col md:flex-row gap-4 mb-6">
            <div className="flex-1 relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
              <Input
                placeholder="Search rules..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-10"
              />
            </div>
            <Select value={actionFilter} onValueChange={setActionFilter}>
              <SelectTrigger className="w-full md:w-[180px]">
                <SelectValue placeholder="Action" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Actions</SelectItem>
                <SelectItem value="allow">Allow</SelectItem>
                <SelectItem value="deny">Deny</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="border rounded-lg">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Status</TableHead>
                  <TableHead>Priority</TableHead>
                  <TableHead>Name</TableHead>
                  <TableHead>Source</TableHead>
                  <TableHead>Destination</TableHead>
                  <TableHead>Port</TableHead>
                  <TableHead>Protocol</TableHead>
                  <TableHead>Action</TableHead>
                  <TableHead>Hits</TableHead>
                  <TableHead>Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filteredRules.map((rule) => (
                  <TableRow key={rule.id}>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        {getStatusIcon(rule.status)}
                      </div>
                    </TableCell>
                    <TableCell>
                      <Badge variant="outline">{rule.priority}</Badge>
                    </TableCell>
                    <TableCell>
                      <div>
                        <p className="font-medium">{rule.name}</p>
                        <p className="text-xs text-muted-foreground">{rule.id}</p>
                      </div>
                    </TableCell>
                    <TableCell className="font-mono text-sm">{rule.source}</TableCell>
                    <TableCell className="font-mono text-sm">{rule.destination}</TableCell>
                    <TableCell className="font-mono text-sm">{rule.port}</TableCell>
                    <TableCell>{rule.protocol}</TableCell>
                    <TableCell>
                      <Badge variant={getActionColor(rule.action)}>
                        {rule.action}
                      </Badge>
                    </TableCell>
                    <TableCell>{rule.hits.toLocaleString()}</TableCell>
                    <TableCell>
                      <div className="flex gap-2">
                        <Button variant="ghost" size="sm">
                          <Edit className="w-4 h-4" />
                        </Button>
                        <Button variant="ghost" size="sm">
                          <Trash2 className="w-4 h-4" />
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

      {/* Rule Statistics */}
      <div className="grid gap-4 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Top Rules by Hits</CardTitle>
            <CardDescription>Most frequently matched rules</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {rules
                .sort((a, b) => b.hits - a.hits)
                .slice(0, 5)
                .map((rule, index) => (
                  <div key={index} className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <Shield className="w-4 h-4 text-primary" />
                      <div>
                        <p className="text-sm font-medium">{rule.name}</p>
                        <p className="text-xs text-muted-foreground">{rule.id}</p>
                      </div>
                    </div>
                    <Badge variant="outline">{rule.hits.toLocaleString()} hits</Badge>
                  </div>
                ))}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Rule Distribution</CardTitle>
            <CardDescription>Rules by action type</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium">Allow Rules</span>
                <div className="flex items-center gap-2">
                  <div className="w-32 h-2 bg-muted rounded-full overflow-hidden">
                    <div 
                      className="h-full bg-green-500" 
                      style={{ width: `${(allowRules / rules.length) * 100}%` }}
                    />
                  </div>
                  <span className="text-sm text-muted-foreground w-12 text-right">
                    {allowRules}
                  </span>
                </div>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium">Deny Rules</span>
                <div className="flex items-center gap-2">
                  <div className="w-32 h-2 bg-muted rounded-full overflow-hidden">
                    <div 
                      className="h-full bg-red-500" 
                      style={{ width: `${(denyRules / rules.length) * 100}%` }}
                    />
                  </div>
                  <span className="text-sm text-muted-foreground w-12 text-right">
                    {denyRules}
                  </span>
                </div>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium">Enabled Rules</span>
                <div className="flex items-center gap-2">
                  <div className="w-32 h-2 bg-muted rounded-full overflow-hidden">
                    <div 
                      className="h-full bg-blue-500" 
                      style={{ width: `${(enabledRules / rules.length) * 100}%` }}
                    />
                  </div>
                  <span className="text-sm text-muted-foreground w-12 text-right">
                    {enabledRules}
                  </span>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
