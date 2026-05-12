const API_BASE_URL = '/api';

export interface Track {
  id: string;
  title: string;
  duration: number; // in nanoseconds from Go, we might need to adjust this
}

export interface Station {
  id: string;
  name: string;
  description: string;
  playlist: Track[];
  ownerId?: number;
}

export interface StationState {
  stationId: string;
  currentTrack: Track;
  startTime: number;
  currentSeconds: number;
}

export async function fetchStations(): Promise<Station[]> {
  const token = localStorage.getItem('token');
  const headers: HeadersInit = {};
  if (token) headers['Authorization'] = `Bearer ${token}`;

  const response = await fetch(`${API_BASE_URL}/stations`, { headers });
  if (!response.ok) throw new Error('Failed to fetch stations');
  return response.json();
}

export async function fetchStationState(id: string): Promise<StationState> {
  const response = await fetch(`${API_BASE_URL}/stations/${id}/state`);
  if (!response.ok) throw new Error('Failed to fetch station state');
  return response.json();
}

export async function createStation(data: Partial<Station>): Promise<Station> {
  const token = localStorage.getItem('token');
  const headers: HeadersInit = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;

  const response = await fetch(`${API_BASE_URL}/stations`, {
    method: 'POST',
    headers,
    body: JSON.stringify(data),
  });
  if (!response.ok) throw new Error('Failed to create station');
  return response.json();
}

export async function updateStation(id: string, data: Partial<Station>): Promise<Station> {
  const token = localStorage.getItem('token');
  const headers: HeadersInit = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;

  const response = await fetch(`${API_BASE_URL}/stations/${id}`, {
    method: 'PUT',
    headers,
    body: JSON.stringify(data),
  });
  if (!response.ok) throw new Error('Failed to update station');
  return response.json();
}

export async function deleteStation(id: string): Promise<void> {
  const token = localStorage.getItem('token');
  const headers: HeadersInit = {};
  if (token) headers['Authorization'] = `Bearer ${token}`;

  const response = await fetch(`${API_BASE_URL}/stations/${id}`, {
    method: 'DELETE',
    headers
  });
  if (!response.ok) throw new Error('Failed to delete station');
}

export async function shareStation(id: string, email: string): Promise<void> {
  const token = localStorage.getItem('token');
  const headers: HeadersInit = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;

  const response = await fetch(`${API_BASE_URL}/stations/${id}/share`, {
    method: 'POST',
    headers,
    body: JSON.stringify({ userEmail: email }),
  });
  if (!response.ok) throw new Error('Failed to share station');
}

export interface PlaylistShare {
  id: number;
  stationId: string;
  stationName?: string;
  userId: number;
  status: string;
}

export async function fetchInvitations(): Promise<PlaylistShare[]> {
  const token = localStorage.getItem('token');
  const headers: HeadersInit = {};
  if (token) headers['Authorization'] = `Bearer ${token}`;

  const response = await fetch(`${API_BASE_URL}/auth/invitations`, { headers });
  if (!response.ok) throw new Error('Failed to fetch invitations');
  return response.json();
}

export async function updateInvitation(id: number, status: 'accepted' | 'rejected'): Promise<void> {
  const token = localStorage.getItem('token');
  const headers: HeadersInit = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;

  const response = await fetch(`${API_BASE_URL}/auth/invitations/${id}`, {
    method: 'PUT',
    headers,
    body: JSON.stringify({ status }),
  });
  if (!response.ok) throw new Error('Failed to update invitation');
}

export interface User {
  id: number;
  username: string;
  email: string;
  isPremium: boolean;
  createdAt: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export async function register(data: Partial<User> & { password?: string }): Promise<AuthResponse> {
  const response = await fetch(`${API_BASE_URL}/auth/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  if (!response.ok) throw new Error('Registration failed');
  return response.json();
}

export async function login(data: { email?: string; password?: string }): Promise<AuthResponse> {
  const response = await fetch(`${API_BASE_URL}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  if (!response.ok) throw new Error('Login failed');
  return response.json();
}
