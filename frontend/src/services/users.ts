import { API_ENDPOINTS } from "../constants";
import { apiClient } from "./apiClient";
import type { UsersResponse } from "../types";

export const fetchUsers = async (): Promise<UsersResponse> => {
	const response = await apiClient.get<UsersResponse>(API_ENDPOINTS.USERS.BASE);
	return response.data;
};
