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
import type { MouseEvent } from "react";
import { formatDate } from "../../utils";
import type { GroupItem as GroupItemType } from "../../types/entities";

interface GroupItemProps {
  group: GroupItemType;
  isLast: boolean;
  onClick: (groupUUID: string) => void;
  createdAtLabel: string;
  onEdit?: (group: GroupItemType) => void;
  onDelete?: (group: GroupItemType) => void;
  editLabel?: string;
  deleteLabel?: string;
}

const GroupItem = ({
  group,
  isLast,
  onClick,
  createdAtLabel,
  onEdit,
  onDelete,
  editLabel,
  deleteLabel,
}: GroupItemProps) => {
  const hasActions = Boolean(onEdit || onDelete);

  const handleEditClick = (event: MouseEvent<HTMLButtonElement>) => {
    event.stopPropagation();
    onEdit?.(group);
  };

  const handleDeleteClick = (event: MouseEvent<HTMLButtonElement>) => {
    event.stopPropagation();
    onDelete?.(group);
  };

  return (
    <Box>
      <ListItem disablePadding alignItems="flex-start">
        <Stack
          spacing={1.5}
          direction={{ xs: "column", sm: "row" }}
          sx={{ width: "100%", alignItems: { xs: "stretch", sm: "center" } }}
        >
          <ListItemButton
            sx={{ borderRadius: 2, flex: 1 }}
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

          {hasActions && (
            <Stack
              direction="row"
              spacing={1}
              sx={{ justifyContent: "flex-end" }}
            >
              {onEdit && (
                <Button size="small" onClick={handleEditClick}>
                  {editLabel}
                </Button>
              )}
              {onDelete && (
                <Button size="small" color="error" onClick={handleDeleteClick}>
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
