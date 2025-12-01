import { Box, CircularProgress } from "@mui/material";
import type React from "react";
import { Navigate, Outlet, useLocation } from "react-router-dom";
import { ROUTES } from "../constants";
import { useAuth } from "../contexts/AuthContextBase";

interface ProtectedRouteProps {
	children?: React.ReactNode;
}

export const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children }) => {
	const { isAuthenticated, isLoading } = useAuth();
	const location = useLocation();

	// Show loading spinner while checking authentication
	if (isLoading) {
		return (
			<Box
				display="flex"
				justifyContent="center"
				alignItems="center"
				minHeight="100vh"
			>
				<CircularProgress />
			</Box>
		);
	}

	// If not authenticated, redirect to login page
	if (!isAuthenticated) {
		return <Navigate to={ROUTES.LOGIN} state={{ from: location }} replace />;
	}

	// If authenticated, render the protected content
	if (children) {
		return <>{children}</>;
	}

	return <Outlet />;
};
