import { Box, Stack } from "@mui/material";
import { useNavigate } from "react-router-dom";
import { ActionButton } from "../components/forms";
import { ROUTES } from "../constants";
import { useLanguage } from "../contexts/LanguageContext";
import type { BaseComponentProps } from "../types";

type HomeProps = BaseComponentProps;

const Home = ({ className = "" }: HomeProps) => {
	const { t } = useLanguage();
	const navigate = useNavigate();

	return (
		<Box
			className={className}
			sx={{
				flex: 1,
				display: "flex",
				alignItems: "center",
				justifyContent: "center",
				minHeight: { xs: "60vh", md: "100%" },
				px: { xs: 2, md: 4 },
				py: { xs: 3, md: 4 },
			}}
		>
			<Stack
				direction={{ xs: "column", sm: "row" }}
				spacing={3}
				justifyContent="center"
				alignItems="center"
			>
				<ActionButton onClick={() => navigate(ROUTES.DOCUMENTS)}>
					{t("home.createDocument")}
				</ActionButton>
				<ActionButton onClick={() => navigate(ROUTES.GROUPS)}>
					{t("home.createGroup")}
				</ActionButton>
			</Stack>
		</Box>
	);
};

export default Home;
