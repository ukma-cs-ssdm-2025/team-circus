import DeleteOutlineIcon from "@mui/icons-material/DeleteOutline";
import {
	Box,
	Card,
	CardActionArea,
	CardContent,
	Chip,
	IconButton,
	Stack,
	Tooltip,
	Typography,
} from "@mui/material";
import { Link as RouterLink } from "react-router-dom";
import { ROUTES } from "../../constants";
import type { DocumentItem } from "../../types/entities";
import { formatDate, truncateText } from "../../utils";

interface DocumentCardProps {
	document: DocumentItem;
	groupName: string;
	createdAtLabel: string;
	noContentLabel: string;
	groupUnknownLabel: string;
	deleteLabel?: string;
	onDelete?: (document: DocumentItem) => void;
}

const DocumentCard = ({
	document,
	groupName,
	createdAtLabel,
	noContentLabel,
	groupUnknownLabel,
	deleteLabel,
	onDelete,
}: DocumentCardProps) => {
	const documentPath = `${ROUTES.DOCUMENTS}/${document.uuid}`;
	const handleDeleteClick = (event: React.MouseEvent<HTMLButtonElement>) => {
		event.preventDefault();
		event.stopPropagation();
		onDelete?.(document);
	};

	return (
		<Card
			sx={{
				height: "100%",
				display: "flex",
				flexDirection: "column",
				borderRadius: 3,
				boxShadow: "0 10px 30px rgba(0, 0, 0, 0.08)",
				transition: "transform 0.2s ease, box-shadow 0.2s ease",
				"&:hover": {
					transform: "translateY(-4px)",
					boxShadow: "0 16px 32px rgba(0, 0, 0, 0.12)",
				},
			}}
		>
			<CardActionArea
				component={RouterLink}
				to={documentPath}
				sx={{ height: "100%", alignItems: "stretch", display: "flex" }}
			>
				<CardContent sx={{ flex: 1 }}>
					<Stack spacing={2}>
						<Stack
							direction="row"
							spacing={1}
							alignItems="center"
							sx={{ minWidth: 0, flexWrap: "wrap", rowGap: 0.5 }}
						>
							<Typography
								variant="h6"
								sx={{
									fontWeight: 600,
									minWidth: 0,
									flexShrink: 1,
									wordBreak: "break-word",
									overflowWrap: "anywhere",
									display: "-webkit-box",
									WebkitLineClamp: 2,
									WebkitBoxOrient: "vertical",
									overflow: "hidden",
								}}
							>
								{document.name}
							</Typography>
							<Chip
								label={groupName || groupUnknownLabel}
								size="small"
								color="primary"
								sx={{ fontWeight: 600 }}
							/>
							{onDelete && (
								<Box ml="auto">
									<Tooltip title={deleteLabel || ""}>
										<IconButton
											color="error"
											size="small"
											onClick={handleDeleteClick}
										>
											<DeleteOutlineIcon fontSize="small" />
										</IconButton>
									</Tooltip>
								</Box>
							)}
						</Stack>

						<Typography
							variant="body2"
							color="text.secondary"
							sx={{ overflowWrap: "anywhere", wordBreak: "break-word" }}
						>
							{truncateText(document.content || noContentLabel, 180)}
						</Typography>

						<Typography variant="body2" color="text.secondary">
							{`${createdAtLabel}: ${formatDate(document.created_at)}`}
						</Typography>
					</Stack>
				</CardContent>
			</CardActionArea>
		</Card>
	);
};

export default DocumentCard;
