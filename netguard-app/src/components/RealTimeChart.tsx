'use client';

import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { useEffect, useState } from 'react';

interface RealTimeChartProps {
  stats: {
    packetsPerSecond: number;
    bytesPerSecond: number;
    activeConnections: number;
    threatsDetected: number;
  };
}

export default function RealTimeChart({ stats }: RealTimeChartProps) {
  const [data, setData] = useState<any[]>([]);

  useEffect(() => {
    const newPoint = {
      time: new Date().toLocaleTimeString(),
      packets: stats.packetsPerSecond,
      bandwidth: (stats.bytesPerSecond / 1024 / 1024).toFixed(2),
    };

    setData((prev) => [...prev.slice(-19), newPoint]);
  }, [stats]);

  return (
    <div className="bg-gray-800/50 backdrop-blur-sm rounded-lg p-6 border border-gray-700">
      <h2 className="text-xl font-bold text-white mb-4">Real-Time Traffic</h2>
      <ResponsiveContainer width="100%" height={300}>
        <LineChart data={data}>
          <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
          <XAxis dataKey="time" stroke="#9CA3AF" />
          <YAxis stroke="#9CA3AF" />
          <Tooltip
            contentStyle={{
              backgroundColor: '#1F2937',
              border: '1px solid #374151',
              borderRadius: '0.5rem',
            }}
          />
          <Line type="monotone" dataKey="packets" stroke="#3B82F6" strokeWidth={2} dot={false} />
          <Line type="monotone" dataKey="bandwidth" stroke="#10B981" strokeWidth={2} dot={false} />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}
