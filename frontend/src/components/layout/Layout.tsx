import { Box, useTheme } from "@mui/material";
import { useState } from "react";
import type { BaseComponentProps } from "../../types";
import Footer from "./Footer";
import Header from "./Header";
import Sidebar from "./Sidebar";

interface LayoutProps extends BaseComponentProps {
	children: React.ReactNode;
	hideFooter?: boolean;
	noPadding?: boolean;
}

const Layout = ({
	children,
	className = "",
	hideFooter = false,
	noPadding = false,
}: LayoutProps) => {
	const theme = useTheme();
	const [isSidebarOpen, setIsSidebarOpen] = useState(false);

	const background =
		noPadding || hideFooter
			? theme.palette.background.default
			: theme.palette.mode === "light"
				? "linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%)"
				: "linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%)";

	const handleOpenSidebar = () => setIsSidebarOpen(true);
	const handleCloseSidebar = () => setIsSidebarOpen(false);

	return (
		<Box
			className={className}
			sx={{
				minHeight: "100vh",
				background,
				display: "flex",
				flexDirection: "column",
			}}
		>
			<Header onToggleSidebar={handleOpenSidebar} />
			<Box
				component="main"
				sx={{
					flex: 1,
					display: "flex",
					flexDirection: "column",
					px: noPadding ? 0 : { xs: 1.5, md: 3 },
					py: noPadding ? 0 : { xs: 2, md: 3 },
				}}
			>
				<Box
					component="section"
					sx={{
						flex: 1,
						display: "flex",
						flexDirection: "column",
						width: "100%",
					}}
				>
					{children}
				</Box>
			</Box>
			{!hideFooter && <Footer />}
			<Sidebar open={isSidebarOpen} onClose={handleCloseSidebar} />
		</Box>
	);
};

export default Layout;
