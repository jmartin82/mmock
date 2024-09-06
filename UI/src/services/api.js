import axios from 'axios';

const API_URL = window.location.protocol + "//" + location.host + "/api";

export const getMockDefinitions = async () => {
  try {
    const response = await axios.get(`${API_URL}/mapping`);
    return response.data;
  } catch (error) {
    console.error('Failed to fetch mock definitions:', error);
    throw error;
  }
};

export const getMockDefinition = async (mockPath) => {
  try {
    const response = await axios.get(`${API_URL}/mapping/${mockPath}`);
    return response.data;
  } catch (error) {
    console.error(`Failed to fetch mock definition for ${mockPath}:`, error);
    throw error;
  }
};

export const createMockDefinition = async (mockPath, mockDefinition) => {
  try {
    const response = await axios.post(`${API_URL}/mapping/${mockPath}`, mockDefinition);
    return response.data;
  } catch (error) {
    console.error('Failed to create mock definition:', error);
    throw error;
  }
};

export const updateMockDefinition = async (mockPath, mockDefinition) => {
  try {
    const response = await axios.put(`${API_URL}/mapping/${mockPath}`, mockDefinition);
    return response.data;
  } catch (error) {
    console.error(`Failed to update mock definition for ${mockPath}:`, error);
    throw error;
  }
};

export const deleteMockDefinition = async (mockPath) => {
  try {
    const response = await axios.delete(`${API_URL}/mapping/${mockPath}`);
    return response.data;
  } catch (error) {
    console.error(`Failed to delete mock definition for ${mockPath}:`, error);
    throw error;
  }
};