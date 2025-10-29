import { List, Stack, Typography } from "@mui/material";
import GroupItem from "./GroupItem";
import type {
  GroupItem as GroupItemType,
  GroupRole,
} from "../../types/entities";

interface GroupsListProps {
  groups: GroupItemType[];
  onGroupClick: (groupUUID: string) => void;
  totalLabel: string;
  createdAtLabel: string;
  onGroupEdit?: (group: GroupItemType) => void;
  onGroupDelete?: (group: GroupItemType) => void;
  editLabel?: string;
  deleteLabel?: string;
  onGroupManageMembers?: (group: GroupItemType) => void;
  manageMembersLabel?: string;
  roleLabel?: string;
  roleNames?: Partial<Record<GroupRole, string>>;
  canManageMembers?: (group: GroupItemType) => boolean;
  canEditGroup?: (group: GroupItemType) => boolean;
  canDeleteGroup?: (group: GroupItemType) => boolean;
}

const GroupsList = ({
  groups,
  onGroupClick,
  totalLabel,
  createdAtLabel,
  onGroupEdit,
  onGroupDelete,
  editLabel,
  deleteLabel,
  onGroupManageMembers,
  manageMembersLabel,
  roleLabel,
  roleNames,
  canManageMembers,
  canEditGroup,
  canDeleteGroup,
}: GroupsListProps) => {
  return (
    <Stack spacing={2}>
      <Typography variant="body2" color="text.secondary">
        {`${totalLabel} ${groups.length}`}
      </Typography>

      <List disablePadding>
        {groups.map((group, index) => {
          const allowEdit = canEditGroup ? canEditGroup(group) : true;
          const allowDelete = canDeleteGroup ? canDeleteGroup(group) : true;

          return (
            <GroupItem
              key={group.uuid}
              group={group}
              isLast={index === groups.length - 1}
              onClick={onGroupClick}
              createdAtLabel={createdAtLabel}
              onEdit={allowEdit ? onGroupEdit : undefined}
              onDelete={allowDelete ? onGroupDelete : undefined}
              editLabel={editLabel}
              deleteLabel={deleteLabel}
              onManageMembers={canManageMembers?.(group) ? onGroupManageMembers : undefined}
              manageMembersLabel={manageMembersLabel}
              roleLabel={roleLabel}
              roleNames={roleNames}
            />
          );
        })}
      </List>
    </Stack>
  );
};

export default GroupsList;
