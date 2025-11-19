import React, { useEffect, useRef } from "react";
import { Navigate, Outlet, useLocation } from "react-router-dom";
import { useAuth } from "../contexts/AuthContextBase";
import { ROUTES } from "../constants";
import { Box, CircularProgress } from "@mui/material";

interface ProtectedRouteProps {
	children?: React.ReactNode;
}

export const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children }) => {
	const { isAuthenticated, isLoading, refreshToken } = useAuth();
	const location = useLocation();
	const didAttemptRefreshRef = useRef(false);

	useEffect(() => {
		// If not authenticated and not loading, try to refresh token
		if (!isAuthenticated && !isLoading && !didAttemptRefreshRef.current) {
			didAttemptRefreshRef.current = true;
			refreshToken();
		}
	}, [isAuthenticated, isLoading, refreshToken]);

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
