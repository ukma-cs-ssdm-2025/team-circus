import React from "react";
import {
	Box,
	Card,
	CardContent,
	Typography,
	Button,
	Alert,
	Chip,
	Divider,
} from "@mui/material";
import { useAuth } from "../contexts/AuthContextBase";

export const AuthTest: React.FC = () => {
	const { user, isAuthenticated, isLoading, logout, refreshToken } = useAuth();

	const handleRefreshToken = async () => {
		try {
			const success = await refreshToken();
			if (success) {
				alert("Token refreshed successfully!");
			} else {
				alert("Token refresh failed!");
			}
		} catch (error) {
			alert(`Token refresh error: ${error}`);
		}
	};

	const handleLogout = async () => {
		try {
			await logout();
			alert("Logged out successfully!");
		} catch (error) {
			alert(`Logout error: ${error}`);
		}
	};

	if (isLoading) {
		return (
			<Box
				display="flex"
				justifyContent="center"
				alignItems="center"
				minHeight="50vh"
			>
				<Typography>Loading authentication status...</Typography>
			</Box>
		);
	}

	return (
		<Box maxWidth="600px" mx="auto" p={2}>
			<Card>
				<CardContent>
					<Typography variant="h5" gutterBottom>
						Authentication Status
					</Typography>

					<Box mb={2}>
						<Typography variant="body1" gutterBottom>
							<strong>Status:</strong>{" "}
							<Chip
								label={isAuthenticated ? "Authenticated" : "Not Authenticated"}
								color={isAuthenticated ? "success" : "error"}
								size="small"
							/>
						</Typography>
					</Box>

					{isAuthenticated && user ? (
						<Box>
							<Typography variant="h6" gutterBottom>
								User Information
							</Typography>
							<Box mb={2}>
								<Typography variant="body2">
									<strong>UUID:</strong> {user.uuid}
								</Typography>
								<Typography variant="body2">
									<strong>Login:</strong> {user.login}
								</Typography>
								<Typography variant="body2">
									<strong>Email:</strong> {user.email}
								</Typography>
								<Typography variant="body2">
									<strong>Created At:</strong>{" "}
									{new Date(user.createdAt).toLocaleString()}
								</Typography>
							</Box>

							<Divider sx={{ my: 2 }} />

							<Box display="flex" gap={2} flexWrap="wrap">
								<Button
									variant="outlined"
									onClick={handleRefreshToken}
									color="primary"
								>
									Test Token Refresh
								</Button>
								<Button
									variant="contained"
									onClick={handleLogout}
									color="error"
								>
									Logout
								</Button>
							</Box>
						</Box>
					) : (
						<Alert severity="info">
							You are not authenticated. Please login to access protected
							content.
						</Alert>
					)}
				</CardContent>
			</Card>
		</Box>
	);
};
