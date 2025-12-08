import { Card, Stack, Typography } from "@mui/material";
import type { BaseComponentProps } from "../../types";

interface SettingsSectionProps extends BaseComponentProps {
	title: string;
	children: React.ReactNode;
	padding?: number;
}

const SettingsSection = ({
	title,
	children,
	padding = 3,
	className = "",
}: SettingsSectionProps) => {
	return (
		<Card variant="outlined" sx={{ p: padding }} className={className}>
			<Typography variant="h5" gutterBottom>
				{title}
			</Typography>
			<Stack spacing={3}>{children}</Stack>
		</Card>
	);
};

export default SettingsSection;
