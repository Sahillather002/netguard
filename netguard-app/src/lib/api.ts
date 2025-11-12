import axios from 'axios';

const RUST_API_URL = process.env.NEXT_PUBLIC_RUST_API_URL || 'http://localhost:8085';
const GO_GATEWAY_URL = process.env.NEXT_PUBLIC_GO_GATEWAY_URL || 'http://localhost:8080';

export const api = axios.create({
  baseURL: RUST_API_URL,
  timeout: 10000,
});

export const gatewayApi = axios.create({
  baseURL: GO_GATEWAY_URL,
  timeout: 10000,
});

// Add auth token to requests
const addAuthToken = (config: any) => {
  const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null;
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
};

api.interceptors.request.use(addAuthToken);
gatewayApi.interceptors.request.use(addAuthToken);

// Rust API endpoints
export const rustApi = {
  // Health & Metrics
  health: () => api.get('/health'),
  metrics: () => api.get('/metrics'),
  
  // Stats
  getRealtimeStats: () => api.get('/api/v1/stats/realtime'),
  
  // Threats
  detectThreat: (data: {
    source_ip: string;
    destination_ip: string;
    port: number;
    protocol: string;
  }) => api.post('/api/v1/threats/detect', data),
  
  // Capture
  startCapture: (data: { interface: string; filter?: string }) => 
    api.post('/api/v1/capture/start', data),
  stopCapture: () => api.post('/api/v1/capture/stop'),
  getCaptureStatus: () => api.get('/api/v1/capture/status'),
  
  // Firewall
  getFirewallRules: () => api.get('/api/v1/firewall/rules'),
  addFirewallRule: (data: {
    action: string;
    source_ip: string;
    destination_port: number;
    protocol: string;
  }) => api.post('/api/v1/firewall/rules', data),
  deleteFirewallRule: (id: string) => api.delete(`/api/v1/firewall/rules/${id}`),
  
  // Interfaces
  listInterfaces: () => api.get('/api/v1/interfaces'),
};

// Gateway API endpoints
export const gateway = {
  health: () => gatewayApi.get('/health'),
  
  // Dashboard
  getDashboardStats: () => gatewayApi.get('/api/v1/dashboard/stats'),
  getRecentActivity: () => gatewayApi.get('/api/v1/dashboard/recent-activity'),
  
  // Alerts
  listAlerts: (params?: { page?: number; limit?: number; severity?: string }) => 
    gatewayApi.get('/api/v1/alerts', { params }),
  getAlert: (id: string) => gatewayApi.get(`/api/v1/alerts/${id}`),
  createAlert: (data: any) => gatewayApi.post('/api/v1/alerts', data),
  updateAlert: (id: string, data: any) => gatewayApi.put(`/api/v1/alerts/${id}`, data),
  deleteAlert: (id: string) => gatewayApi.delete(`/api/v1/alerts/${id}`),
  
  // Network
  listInterfaces: () => gatewayApi.get('/api/v1/network/interfaces'),
  getNetworkStats: () => gatewayApi.get('/api/v1/network/stats'),
  startMonitoring: (data: any) => gatewayApi.post('/api/v1/network/monitor/start', data),
  stopMonitoring: () => gatewayApi.post('/api/v1/network/monitor/stop'),
  
  // Firewall
  listFirewallRules: () => gatewayApi.get('/api/v1/firewall/rules'),
  addFirewallRule: (data: any) => gatewayApi.post('/api/v1/firewall/rules', data),
  deleteFirewallRule: (id: string) => gatewayApi.delete(`/api/v1/firewall/rules/${id}`),
  
  // Threats
  listThreats: () => gatewayApi.get('/api/v1/threats'),
  getThreat: (id: string) => gatewayApi.get(`/api/v1/threats/${id}`),
  analyzeThreat: (data: any) => gatewayApi.post('/api/v1/threats/analyze', data),
  
  // Users
  listUsers: () => gatewayApi.get('/api/v1/users'),
  getUser: (id: string) => gatewayApi.get(`/api/v1/users/${id}`),
  updateUser: (id: string, data: any) => gatewayApi.put(`/api/v1/users/${id}`, data),
  deleteUser: (id: string) => gatewayApi.delete(`/api/v1/users/${id}`),
};

// WebSocket URLs
export const WS_URLS = {
  threats: process.env.NEXT_PUBLIC_WS_THREATS_URL || 'ws://localhost:8080/ws/threats',
  stats: process.env.NEXT_PUBLIC_WS_STATS_URL || 'ws://localhost:8080/ws/stats',
};
