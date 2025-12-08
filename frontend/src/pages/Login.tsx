import {
	Alert,
	Box,
	Button,
	Card,
	CardContent,
	Container,
	Link,
	Paper,
	TextField,
	Typography,
} from "@mui/material";
import type React from "react";
import { useState } from "react";
import { Navigate, Link as RouterLink, useLocation } from "react-router-dom";
import { ROUTES } from "../constants";
import { useAuth } from "../contexts/AuthContextBase";
import type { LoginRequest } from "../types/auth";

type AuthLocationState = { from?: { pathname?: string } };

const resolveRedirectPath = (state: unknown): string | undefined => {
	if (typeof state !== "object" || state === null) {
		return undefined;
	}

	const maybeState = state as AuthLocationState;
	const pathname = maybeState.from?.pathname;
	return typeof pathname === "string" ? pathname : undefined;
};

export const Login: React.FC = () => {
	const { login, isAuthenticated, isLoading } = useAuth();
	const location = useLocation();
	const [formData, setFormData] = useState<LoginRequest>({
		login: "",
		password: "",
	});
	const [error, setError] = useState<string>("");
	const [isSubmitting, setIsSubmitting] = useState(false);

	// Redirect if already authenticated
	if (isAuthenticated) {
		const from = resolveRedirectPath(location.state) ?? ROUTES.HOME;
		return <Navigate to={from} replace />;
	}

	const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
		const { name, value } = e.target;
		setFormData((prev) => ({
			...prev,
			[name]: value,
		}));
		// Clear error when user starts typing
		if (error) setError("");
	};

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault();
		setIsSubmitting(true);
		setError("");

		try {
			await login(formData);
		} catch (err) {
			setError(err instanceof Error ? err.message : "Login failed");
		} finally {
			setIsSubmitting(false);
		}
	};

	if (isLoading) {
		return (
			<Container maxWidth="sm">
				<Box
					display="flex"
					justifyContent="center"
					alignItems="center"
					minHeight="100vh"
				>
					<Typography>Loading...</Typography>
				</Box>
			</Container>
		);
	}

	return (
		<Container maxWidth="sm">
			<Box
				display="flex"
				justifyContent="center"
				alignItems="center"
				minHeight="100vh"
				py={4}
			>
				<Paper elevation={3} sx={{ width: "100%", maxWidth: 400 }}>
					<Card>
						<CardContent sx={{ p: 4 }}>
							<Typography
								variant="h4"
								component="h1"
								gutterBottom
								align="center"
							>
								Sign In
							</Typography>
							<Typography
								variant="body2"
								color="text.secondary"
								align="center"
								sx={{ mb: 3 }}
							>
								Welcome back! Please sign in to your account.
							</Typography>

							{error && (
								<Alert severity="error" sx={{ mb: 2 }}>
									{error}
								</Alert>
							)}

							<Box component="form" onSubmit={handleSubmit}>
								<TextField
									fullWidth
									label="Username"
									name="login"
									value={formData.login}
									onChange={handleChange}
									margin="normal"
									required
									autoComplete="username"
									autoFocus
								/>
								<TextField
									fullWidth
									label="Password"
									name="password"
									type="password"
									value={formData.password}
									onChange={handleChange}
									margin="normal"
									required
									autoComplete="current-password"
								/>
								<Button
									type="submit"
									fullWidth
									variant="contained"
									sx={{ mt: 3, mb: 2 }}
									disabled={isSubmitting}
								>
									{isSubmitting ? "Signing In..." : "Sign In"}
								</Button>
								<Box textAlign="center">
									<Typography variant="body2">
										Don't have an account?{" "}
										<Link component={RouterLink} to={ROUTES.REGISTER}>
											Sign up here
										</Link>
									</Typography>
								</Box>
							</Box>
						</CardContent>
					</Card>
				</Paper>
			</Box>
		</Container>
	);
};
