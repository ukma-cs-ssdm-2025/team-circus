import {
	Box,
	Button,
	Divider,
	ListItem,
	ListItemButton,
	ListItemText,
	Stack,
	Typography,
} from "@mui/material";
import { formatDate } from "../../utils";
import type { GroupItem as GroupItemType } from "../../types/entities";

interface GroupItemProps {
	group: GroupItemType;
	isLast: boolean;
	onClick: (groupUUID: string) => void;
	createdAtLabel: string;
	onManage?: (groupUUID: string) => void;
	onDelete?: (group: GroupItemType) => void;
	manageLabel?: string;
	deleteLabel?: string;
}

const GroupItem = ({
	group,
	isLast,
	onClick,
	createdAtLabel,
	onManage,
	onDelete,
	manageLabel,
	deleteLabel,
}: GroupItemProps) => {
	const handleManage = (event: React.MouseEvent<HTMLButtonElement>) => {
		event.stopPropagation();
		onManage?.(group.uuid);
	};

	const handleDelete = (event: React.MouseEvent<HTMLButtonElement>) => {
		event.stopPropagation();
		onDelete?.(group);
	};

	return (
		<Box>
			<ListItem disablePadding alignItems="flex-start">
				<Stack spacing={1} sx={{ width: "100%" }}>
					<ListItemButton
						sx={{ borderRadius: 2 }}
						onClick={() => onClick(group.uuid)}
					>
						<ListItemText
							primary={
								<Typography variant="h6" sx={{ fontWeight: 600 }}>
									{group.name}
								</Typography>
							}
							secondary={
								<Typography variant="body2" color="text.secondary">
									{`${createdAtLabel}: ${formatDate(group.created_at)}`}
								</Typography>
							}
						/>
					</ListItemButton>

					{(onManage || onDelete) && (
						<Stack
							direction="row"
							spacing={1}
							justifyContent="flex-end"
							px={1}
							pb={1}
						>
							{onManage && (
								<Button size="small" onClick={handleManage}>
									{manageLabel}
								</Button>
							)}
							{onDelete && (
								<Button size="small" color="error" onClick={handleDelete}>
									{deleteLabel}
								</Button>
							)}
						</Stack>
					)}
				</Stack>
			</ListItem>
			{!isLast && <Divider sx={{ my: 1 }} />}
		</Box>
	);
};

export default GroupItem;
