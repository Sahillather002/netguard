'use client'

import { useCallback, useEffect, useState } from 'react';
import ReactFlow, {
  Node,
  Edge,
  addEdge,
  Connection,
  useNodesState,
  useEdgesState,
  Controls,
  Background,
  MiniMap,
  BackgroundVariant,
} from 'reactflow';
import 'reactflow/dist/style.css';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Loader2 } from 'lucide-react';

interface AttackFlowData {
  nodes: Node[];
  edges: Edge[];
}

interface AttackFlowVisualizationProps {
  data?: AttackFlowData;
  isLoading?: boolean;
}

const nodeTypes = {
  // You can define custom node types here
};

const defaultNodes: Node[] = [
  {
    id: '1',
    type: 'input',
    data: { label: 'External Attacker' },
    position: { x: 250, y: 0 },
    style: { background: '#ef4444', color: 'white', border: '2px solid #dc2626' },
  },
  {
    id: '2',
    data: { label: 'Phishing Email' },
    position: { x: 250, y: 100 },
    style: { background: '#f97316', color: 'white' },
  },
  {
    id: '3',
    data: { label: 'User Workstation' },
    position: { x: 100, y: 200 },
    style: { background: '#eab308', color: 'white' },
  },
  {
    id: '4',
    data: { label: 'Malware Execution' },
    position: { x: 400, y: 200 },
    style: { background: '#ef4444', color: 'white' },
  },
  {
    id: '5',
    data: { label: 'Lateral Movement' },
    position: { x: 250, y: 300 },
    style: { background: '#dc2626', color: 'white' },
  },
  {
    id: '6',
    data: { label: 'Database Server' },
    position: { x: 100, y: 400 },
    style: { background: '#7c3aed', color: 'white' },
  },
  {
    id: '7',
    data: { label: 'Data Exfiltration' },
    position: { x: 400, y: 400 },
    style: { background: '#dc2626', color: 'white' },
  },
  {
    id: '8',
    type: 'output',
    data: { label: 'Blocked by Firewall' },
    position: { x: 250, y: 500 },
    style: { background: '#22c55e', color: 'white', border: '2px solid #16a34a' },
  },
];

const defaultEdges: Edge[] = [
  { id: 'e1-2', source: '1', target: '2', animated: true, label: 'Sends' },
  { id: 'e2-3', source: '2', target: '3', animated: true, label: 'Opens' },
  { id: 'e2-4', source: '2', target: '4', animated: true, label: 'Downloads' },
  { id: 'e3-5', source: '3', target: '5', animated: true, label: 'Compromised' },
  { id: 'e4-5', source: '4', target: '5', animated: true, label: 'Spreads' },
  { id: 'e5-6', source: '5', target: '6', animated: true, label: 'Accesses' },
  { id: 'e5-7', source: '5', target: '7', animated: true, label: 'Attempts' },
  { id: 'e7-8', source: '7', target: '8', animated: true, label: 'Detected', style: { stroke: '#22c55e' } },
];

export function AttackFlowVisualization({ 
  data, 
  isLoading = false 
}: AttackFlowVisualizationProps) {
  const [nodes, setNodes, onNodesChange] = useNodesState(data?.nodes || defaultNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(data?.edges || defaultEdges);

  useEffect(() => {
    if (data) {
      setNodes(data.nodes);
      setEdges(data.edges);
    }
  }, [data, setNodes, setEdges]);

  const onConnect = useCallback(
    (params: Connection) => setEdges((eds) => addEdge(params, eds)),
    [setEdges]
  );

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Attack Flow Visualization</CardTitle>
          <CardDescription>Real-time attack path analysis</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="h-[600px] flex items-center justify-center">
            <Loader2 className="w-8 h-8 animate-spin text-primary" />
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle>Attack Flow Visualization</CardTitle>
            <CardDescription>Real-time attack path analysis</CardDescription>
          </div>
          <div className="flex gap-2">
            <Badge variant="destructive">Live</Badge>
            <Badge variant="outline">{nodes.length} Nodes</Badge>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <div className="h-[600px] border rounded-lg overflow-hidden bg-muted/10">
          <ReactFlow
            nodes={nodes}
            edges={edges}
            onNodesChange={onNodesChange}
            onEdgesChange={onEdgesChange}
            onConnect={onConnect}
            nodeTypes={nodeTypes}
            fitView
            attributionPosition="bottom-left"
          >
            <Background variant={BackgroundVariant.Dots} gap={12} size={1} />
            <Controls />
            <MiniMap 
              nodeStrokeColor={(n) => {
                if (n.type === 'input') return '#ef4444';
                if (n.type === 'output') return '#22c55e';
                return '#3b82f6';
              }}
              nodeColor={(n) => {
                return n.style?.background as string || '#3b82f6';
              }}
              nodeBorderRadius={2}
            />
          </ReactFlow>
        </div>
        
        <div className="mt-4 grid grid-cols-2 md:grid-cols-4 gap-4">
          <div className="p-3 border rounded-lg">
            <div className="flex items-center gap-2 mb-1">
              <div className="w-3 h-3 rounded-full bg-red-500" />
              <span className="text-sm font-medium">Attack Source</span>
            </div>
            <p className="text-xs text-muted-foreground">External threats</p>
          </div>
          <div className="p-3 border rounded-lg">
            <div className="flex items-center gap-2 mb-1">
              <div className="w-3 h-3 rounded-full bg-orange-500" />
              <span className="text-sm font-medium">Attack Vector</span>
            </div>
            <p className="text-xs text-muted-foreground">Entry points</p>
          </div>
          <div className="p-3 border rounded-lg">
            <div className="flex items-center gap-2 mb-1">
              <div className="w-3 h-3 rounded-full bg-purple-500" />
              <span className="text-sm font-medium">Target Asset</span>
            </div>
            <p className="text-xs text-muted-foreground">Critical systems</p>
          </div>
          <div className="p-3 border rounded-lg">
            <div className="flex items-center gap-2 mb-1">
              <div className="w-3 h-3 rounded-full bg-green-500" />
              <span className="text-sm font-medium">Blocked</span>
            </div>
            <p className="text-xs text-muted-foreground">Prevented attacks</p>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
