import { CheckCircle, Error, Warning } from "@mui/icons-material";
import { Box, Chip, CircularProgress, Typography } from "@mui/material";
import { ENV } from "../../config";
import { useApi } from "../../hooks";
import type { BaseComponentProps } from "../../types";

interface ApiStatusProps extends BaseComponentProps {
	endpoint?: string;
}

const ApiStatus = ({
	endpoint = "/health",
	className = "",
}: ApiStatusProps) => {
	const { data, loading, error } = useApi<{
		status: string;
		timestamp: string;
	}>(endpoint, {
		immediate: true,
	});

	const getStatusInfo = () => {
		if (loading) {
			return {
				color: "warning" as const,
				icon: <CircularProgress size={16} />,
				label: "Перевірка...",
			};
		}

		if (error) {
			return {
				color: "error" as const,
				icon: <Error />,
				label: "API недоступний",
			};
		}

		if (data?.status === "ok") {
			return {
				color: "success" as const,
				icon: <CheckCircle />,
				label: "API працює",
			};
		}

		return {
			color: "warning" as const,
			icon: <Warning />,
			label: "Невідомий статус",
		};
	};

	const statusInfo = getStatusInfo();

	return (
		<Box
			className={className}
			sx={{ display: "flex", alignItems: "center", gap: 1 }}
		>
			<Typography variant="body2" color="text.secondary">
				API: {ENV.API_BASE_URL}
			</Typography>
			<Chip
				icon={statusInfo.icon}
				label={statusInfo.label}
				color={statusInfo.color}
				size="small"
				variant="outlined"
			/>
		</Box>
	);
};

export default ApiStatus;
