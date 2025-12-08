import { API_ENDPOINTS } from "../constants";
import type { UsersResponse } from "../types";
import { apiClient } from "./apiClient";

export const fetchUsers = async (): Promise<UsersResponse> => {
	const response = await apiClient.get<UsersResponse>(API_ENDPOINTS.USERS.BASE);
	return response.data;
};
