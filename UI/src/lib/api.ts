// API Configuration and Services for MMock Console

const DEFAULT_BASE_URL = window.location.protocol + "//" + window.location.host;

export const getBaseUrl = (): string => {
  return localStorage.getItem('mmock-api-url') || DEFAULT_BASE_URL;
};

export const setBaseUrl = (url: string): void => {
  localStorage.setItem('mmock-api-url', url);
};

export interface MatchError {
  uri: string;
  reason: string;
}

export interface MatchResult {
  match: boolean;
  uri: string;
  errors: MatchError[];
}

export interface Transaction {
  time: number;
  request: {
    scheme: string;
    host: string;
    port: string;
    method: string;
    path: string;
    queryStringParameters: Record<string, string[]>;
    fragment: string;
    headers: Record<string, string[]>;
    cookies: Record<string, string>;
    body: string;
  };
  response: {
    statusCode: number;
    headers: Record<string, string[]> | null;
    cookies: Record<string, string> | null;
    body: string;
  };
  result: MatchResult;
}

export interface MockDefinition {
  URI: string;
  description: string;
  request: {
    scheme?: string;
    host?: string;
    port?: string;
    method: string;
    path: string;
    queryStringParameters?: Record<string, string>;
    headers?: Record<string, string[]>;
    cookies?: Record<string, string>;
    body?: string;
  };
  response: {
    statusCode: number;
    headers?: Record<string, string[]>;
    cookies?: Record<string, string>;
    body?: string;
  };
  control?: {
    priority?: number;
    delay?: string;
    crazy?: boolean;
    proxyBaseURL?: string;
    webHookURL?: string;
    scenario?: {
      name: string;
      requiredState?: string[];
      newState?: string;
    };
  };
  callback?: {
    delay?: string;
    method?: string;
    url?: string;
    body?: string;
    headers?: Record<string, string[]>;
    cookies?: Record<string, string>;
  };
}

export interface Scenarios {
  [scenarioName: string]: string;
}

// API Functions
const fetchApi = async <T>(endpoint: string, options?: RequestInit): Promise<T> => {
  const baseUrl = getBaseUrl();
  const response = await fetch(`${baseUrl}${endpoint}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options?.headers,
    },
  });

  if (!response.ok) {
    throw new Error(`API Error: ${response.status} ${response.statusText}`);
  }

  const text = await response.text();
  return text ? JSON.parse(text) : null;
};

// Request Monitoring
export const requestApi = {
  getAll: () => fetchApi<Transaction[]>('/api/request/all'),
  getPage: (page: number) => fetchApi<Transaction[]>(`/api/request/all/${page}`),
  getMatched: () => fetchApi<Transaction[]>('/api/request/matched'),
  getUnmatched: () => fetchApi<Transaction[]>('/api/request/unmatched'),
  reset: () => fetchApi<void>('/api/request/reset'),
  verify: (request: Partial<Transaction['request']>) =>
    fetchApi<{ matched: boolean }>('/api/request/verify', {
      method: 'POST',
      body: JSON.stringify(request),
    }),
};

// Mock Management
export const mappingApi = {
  getAll: () => fetchApi<MockDefinition[]>('/api/mapping'),
  get: (path: string) => fetchApi<MockDefinition>(`/api/mapping/${encodeURIComponent(path)}`),
  create: (path: string, mock: MockDefinition) =>
    fetchApi<MockDefinition>(`/api/mapping/${encodeURIComponent(path)}`, {
      method: 'POST',
      body: JSON.stringify(mock),
    }),
  update: (path: string, mock: MockDefinition) =>
    fetchApi<MockDefinition>(`/api/mapping/${encodeURIComponent(path)}`, {
      method: 'PUT',
      body: JSON.stringify(mock),
    }),
  delete: (path: string) =>
    fetchApi<void>(`/api/mapping/${encodeURIComponent(path)}`, {
      method: 'DELETE',
    }),
};

// Scenario Management
export const scenarioApi = {
  getAll: () => fetchApi<Scenarios>('/api/scenarios'),
  resetAll: () => fetchApi<void>('/api/scenarios/reset_all'),
  setState: (scenario: string, state: string) =>
    fetchApi<void>(`/api/scenarios/set/${encodeURIComponent(scenario)}/${encodeURIComponent(state)}`, {
      method: 'PUT',
    }),
  pause: () => fetchApi<void>('/api/scenarios/pause', { method: 'PUT' }),
  unpause: () => fetchApi<void>('/api/scenarios/unpause', { method: 'PUT' }),
};
