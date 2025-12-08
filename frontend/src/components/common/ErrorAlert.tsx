import { Alert, Box, Button } from "@mui/material";

interface ErrorAlertProps {
	message: string;
	onRetry?: () => void;
	retryText?: string;
}

const ErrorAlert = ({
	message,
	onRetry,
	retryText = "Retry",
}: ErrorAlertProps) => {
	return (
		<Alert
			severity="error"
			sx={{ mb: 3, display: "flex", alignItems: "center", gap: 2 }}
		>
			<Box component="span">{message}</Box>
			{onRetry && (
				<Button variant="outlined" size="small" onClick={onRetry}>
					{retryText}
				</Button>
			)}
		</Alert>
	);
};

export default ErrorAlert;
