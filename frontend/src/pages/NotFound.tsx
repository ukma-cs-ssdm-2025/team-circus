import {
	ArrowBack as ArrowBackIcon,
	Home as HomeIcon,
} from "@mui/icons-material";
import { Box, Stack, Typography, useTheme } from "@mui/material";
import { Link } from "react-router-dom";
import { CenteredContent, PageCard } from "../components/common";
import { ActionButton } from "../components/forms";
import { ROUTES } from "../constants";
import { useLanguage } from "../contexts/LanguageContext";
import type { BaseComponentProps } from "../types";

type NotFoundProps = BaseComponentProps;

const NotFound = ({ className = "" }: NotFoundProps) => {
	const theme = useTheme();
	const { t } = useLanguage();

	return (
		<CenteredContent className={className} maxWidth="sm">
			<PageCard>
				<Box sx={{ textAlign: "center" }}>
					<Typography
						variant="h1"
						sx={{
							fontSize: "8rem",
							fontWeight: 900,
							color: theme.palette.primary.main,
							mb: 2,
							lineHeight: 1,
						}}
					>
						404
					</Typography>

					<Typography variant="h3" gutterBottom>
						{t("notFound.title")}
					</Typography>

					<Typography variant="h6" color="text.secondary" sx={{ mb: 4 }}>
						{t("notFound.message")}
					</Typography>

					<Stack
						direction={{ xs: "column", sm: "row" }}
						spacing={2}
						justifyContent="center"
					>
						<ActionButton
							component={Link}
							to={ROUTES.HOME}
							startIcon={<HomeIcon />}
						>
							{t("notFound.home")}
						</ActionButton>
						<ActionButton
							variant="outlined"
							startIcon={<ArrowBackIcon />}
							onClick={() => window.history.back()}
						>
							{t("notFound.back")}
						</ActionButton>
					</Stack>
				</Box>
			</PageCard>
		</CenteredContent>
	);
};

export default NotFound;
