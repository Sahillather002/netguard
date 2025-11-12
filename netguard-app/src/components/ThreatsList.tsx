import { AlertTriangle } from 'lucide-react';

interface Threat {
  threat_type?: string;
  severity?: string;
  source_ip?: string;
  timestamp?: string;
}

interface ThreatsListProps {
  threats: Threat[];
}

export default function ThreatsList({ threats }: ThreatsListProps) {
  const getSeverityColor = (severity: string = 'low') => {
    const colors = {
      critical: 'text-red-400 bg-red-500/20',
      high: 'text-orange-400 bg-orange-500/20',
      medium: 'text-yellow-400 bg-yellow-500/20',
      low: 'text-blue-400 bg-blue-500/20',
    };
    return colors[severity.toLowerCase() as keyof typeof colors] || colors.low;
  };

  return (
    <div className="bg-gray-800/50 backdrop-blur-sm rounded-lg p-6 border border-gray-700">
      <h2 className="text-xl font-bold text-white mb-4 flex items-center">
        <AlertTriangle className="w-5 h-5 mr-2 text-red-400" />
        Recent Threats
      </h2>
      <div className="space-y-3 max-h-96 overflow-y-auto">
        {threats.length === 0 ? (
          <p className="text-gray-400 text-center py-8">No threats detected</p>
        ) : (
          threats.map((threat, index) => (
            <div
              key={index}
              className="bg-gray-900/50 rounded-lg p-4 border border-gray-700 hover:border-gray-600 transition-colors"
            >
              <div className="flex items-center justify-between mb-2">
                <span className="text-white font-semibold">
                  {threat.threat_type || 'Unknown Threat'}
                </span>
                <span
                  className={`px-2 py-1 rounded text-xs font-semibold ${
                    getSeverityColor(threat.severity)
                  }`}
                >
                  {threat.severity || 'LOW'}
                </span>
              </div>
              <div className="text-sm text-gray-400">
                <p>Source: {threat.source_ip || 'Unknown'}</p>
                <p className="text-xs mt-1">
                  {threat.timestamp
                    ? new Date(threat.timestamp).toLocaleString()
                    : 'Just now'}
                </p>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
