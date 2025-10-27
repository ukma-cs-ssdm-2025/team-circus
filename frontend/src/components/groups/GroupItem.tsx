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
import type {
  GroupItem as GroupItemType,
  GroupRole,
} from "../../types/entities";

interface GroupItemProps {
  group: GroupItemType;
  isLast: boolean;
  onClick: (groupUUID: string) => void;
  createdAtLabel: string;
  onEdit?: (group: GroupItemType) => void;
  onDelete?: (group: GroupItemType) => void;
  editLabel?: string;
  deleteLabel?: string;
  onManageMembers?: (group: GroupItemType) => void;
  manageMembersLabel?: string;
  roleLabel?: string;
  roleNames?: Partial<Record<GroupRole, string>>;
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
  onManageMembers,
  manageMembersLabel,
  roleLabel,
  roleNames,
}: GroupItemProps) => {
  const hasActions = Boolean(onEdit || onDelete || onManageMembers);

  const handleEditClick = (event: MouseEvent<HTMLButtonElement>) => {
    event.stopPropagation();
    onEdit?.(group);
  };

  const handleDeleteClick = (event: MouseEvent<HTMLButtonElement>) => {
    event.stopPropagation();
    onDelete?.(group);
  };

  const handleManageMembersClick = (event: MouseEvent<HTMLButtonElement>) => {
    event.stopPropagation();
    onManageMembers?.(group);
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
                <Stack direction="row" spacing={1} alignItems="center">
                  <Typography variant="body2" color="text.secondary">
                    {`${createdAtLabel}: ${formatDate(group.created_at)}`}
                  </Typography>
                  {roleLabel && group.role && (
                    <Typography variant="caption" color="primary">
                      {`${roleLabel}: ${roleNames?.[group.role] ?? group.role}`}
                    </Typography>
                  )}
                </Stack>
              }
            />
          </ListItemButton>

          {hasActions && (
            <Stack
              direction="row"
              spacing={1}
              sx={{ justifyContent: "flex-end" }}
            >
              {onManageMembers && (
                <Button size="small" color="info" onClick={handleManageMembersClick}>
                  {manageMembersLabel}
                </Button>
              )}
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
