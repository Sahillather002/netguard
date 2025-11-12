import axios, { AxiosInstance, AxiosRequestConfig } from 'axios';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000/api';

class ApiClient {
  private client: AxiosInstance;

  constructor() {
    this.client = axios.create({
      baseURL: API_URL,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // Request interceptor to add auth token
    this.client.interceptors.request.use(
      (config) => {
        const token = this.getToken();
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
      },
      (error) => Promise.reject(error)
    );

    // Response interceptor for error handling
    this.client.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response?.status === 401) {
          this.clearToken();
          if (typeof window !== 'undefined') {
            window.location.href = '/login';
          }
        }
        return Promise.reject(error);
      }
    );
  }

  private getToken(): string | null {
    if (typeof window === 'undefined') return null;
    return localStorage.getItem('token');
  }

  private setToken(token: string): void {
    if (typeof window !== 'undefined') {
      localStorage.setItem('token', token);
    }
  }

  private clearToken(): void {
    if (typeof window !== 'undefined') {
      localStorage.removeItem('token');
    }
  }

  // Authentication
  async login(email: string, password: string) {
    const response = await this.client.post('/auth/login', { email, password });
    if (response.data.token) {
      this.setToken(response.data.token);
    }
    return response.data;
  }

  async register(data: {
    firstName: string;
    lastName: string;
    email: string;
    password: string;
    company: string;
    role: string;
  }) {
    const response = await this.client.post('/auth/register', data);
    if (response.data.token) {
      this.setToken(response.data.token);
    }
    return response.data;
  }

  async logout() {
    this.clearToken();
    return { success: true };
  }

  async getCurrentUser() {
    const response = await this.client.get('/auth/me');
    return response.data;
  }

  // Dashboard
  async getDashboardStats() {
    const response = await this.client.get('/dashboard/stats');
    return response.data;
  }

  // Alerts
  async getAlerts(params?: { severity?: string; status?: string; search?: string }) {
    const response = await this.client.get('/alerts', { params });
    return response.data;
  }

  async getAlertById(id: string) {
    const response = await this.client.get(`/alerts/${id}`);
    return response.data;
  }

  async updateAlertStatus(id: string, status: string) {
    const response = await this.client.patch(`/alerts/${id}/status`, { status });
    return response.data;
  }

  // Threats
  async getThreats(params?: { type?: string; severity?: string; status?: string }) {
    const response = await this.client.get('/threats', { params });
    return response.data;
  }

  async getThreatById(id: string) {
    const response = await this.client.get(`/threats/${id}`);
    return response.data;
  }

  async blockThreat(id: string) {
    const response = await this.client.post(`/threats/${id}/block`);
    return response.data;
  }

  // Network
  async getNetworkDevices(params?: { search?: string }) {
    const response = await this.client.get('/network/devices', { params });
    return response.data;
  }

  async getNetworkStats() {
    const response = await this.client.get('/network/stats');
    return response.data;
  }

  // Firewall
  async getFirewallRules(params?: { action?: string; search?: string }) {
    const response = await this.client.get('/firewall/rules', { params });
    return response.data;
  }

  async createFirewallRule(data: any) {
    const response = await this.client.post('/firewall/rules', data);
    return response.data;
  }

  async updateFirewallRule(id: string, data: any) {
    const response = await this.client.put(`/firewall/rules/${id}`, data);
    return response.data;
  }

  async deleteFirewallRule(id: string) {
    const response = await this.client.delete(`/firewall/rules/${id}`);
    return response.data;
  }

  // Users
  async getUsers(params?: { role?: string; search?: string }) {
    const response = await this.client.get('/users', { params });
    return response.data;
  }

  async getUserById(id: string) {
    const response = await this.client.get(`/users/${id}`);
    return response.data;
  }

  async createUser(data: any) {
    const response = await this.client.post('/users', data);
    return response.data;
  }

  async updateUser(id: string, data: any) {
    const response = await this.client.put(`/users/${id}`, data);
    return response.data;
  }

  async deleteUser(id: string) {
    const response = await this.client.delete(`/users/${id}`);
    return response.data;
  }

  // Attack Flow Visualization
  async getAttackFlow() {
    const response = await this.client.get('/threats/attack-flow');
    return response.data;
  }
}

export const apiClient = new ApiClient();
export default apiClient;
