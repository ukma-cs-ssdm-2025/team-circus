import { useMemo, useState } from "react";
import {
  Alert,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Stack,
  TextField,
  Typography,
} from "@mui/material";
import { useNavigate } from "react-router-dom";
import {
  CenteredContent,
  PageCard,
  PageHeader,
  ErrorAlert,
  LoadingSpinner,
  GroupsList,
  ConfirmDialog,
} from "../components";
import { useLanguage } from "../contexts/LanguageContext";
import { useApi } from "../hooks";
import { API_ENDPOINTS, ROUTES } from "../constants";
import { createGroup, deleteGroup } from "../services";
import type { BaseComponentProps, GroupItem, GroupsResponse } from "../types";

type GroupsProps = BaseComponentProps;

const Groups = ({ className = "" }: GroupsProps) => {
  const { t } = useLanguage();
  const navigate = useNavigate();
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [groupName, setGroupName] = useState("");
  const [creating, setCreating] = useState(false);
  const [createError, setCreateError] = useState("");
  const [groupToDelete, setGroupToDelete] = useState<GroupItem | null>(null);
  const [deleting, setDeleting] = useState(false);
  const [feedback, setFeedback] = useState<{
    type: "success" | "error";
    message: string;
  } | null>(null);

  const { data, loading, error, refetch } = useApi<GroupsResponse>(
    API_ENDPOINTS.GROUPS.BASE,
  );
  const groups = useMemo(() => data?.groups ?? [], [data]);

  const handleOpenGroupDocuments = (groupUUID: string) => {
    navigate({
      pathname: ROUTES.DOCUMENTS,
      search: `?group=${groupUUID}`,
    });
  };

  const handleManageMembers = (groupUUID: string) => {
    navigate({
      pathname: ROUTES.GROUP_DETAILS.replace(":uuid", groupUUID),
    });
  };

  const handleCreateGroup = async () => {
    if (!groupName.trim()) {
      setCreateError(t("groups.fieldRequired"));
      return;
    }

    setCreating(true);
    try {
      await createGroup({ name: groupName.trim() });
      setFeedback({ type: "success", message: t("groups.createSuccess") });
      setGroupName("");
      setCreateError("");
      setCreateDialogOpen(false);
      await refetch();
    } catch (err) {
      console.error("Failed to create group", err);
      setFeedback({ type: "error", message: t("groups.createError") });
    } finally {
      setCreating(false);
    }
  };

  const confirmDeleteGroup = async () => {
    if (!groupToDelete) {
      return;
    }

    setDeleting(true);
    try {
      await deleteGroup(groupToDelete.uuid);
      setFeedback({ type: "success", message: t("groups.deleteSuccess") });
      setGroupToDelete(null);
      await refetch();
    } catch (err) {
      console.error("Failed to delete group", err);
      setFeedback({ type: "error", message: t("groups.deleteError") });
    } finally {
      setDeleting(false);
    }
  };

  return (
    <CenteredContent className={className}>
      <PageCard>
        <PageHeader title={t("groups.title")} subtitle={t("groups.subtitle")} />

        <Stack spacing={3}>
          <Stack
            direction={{ xs: "column", sm: "row" }}
            justifyContent="flex-end"
          >
            <Button
              variant="contained"
              onClick={() => setCreateDialogOpen(true)}
            >
              {t("groups.createButton")}
            </Button>
          </Stack>

          {feedback && (
            <Alert
              severity={feedback.type}
              onClose={() => setFeedback(null)}
              variant="outlined"
            >
              {feedback.message}
            </Alert>
          )}

          {loading && <LoadingSpinner />}

          {error && (
            <ErrorAlert
              message={t("groups.error")}
              onRetry={refetch}
              retryText={t("groups.refresh")}
            />
          )}

          {!loading && !error && groups.length === 0 && (
            <Typography color="text.secondary" align="center">
              {t("groups.empty")}
            </Typography>
          )}

          {groups.length > 0 && (
            <GroupsList
              groups={groups}
              onGroupClick={handleOpenGroupDocuments}
              onManageMembers={handleManageMembers}
              onDeleteGroup={setGroupToDelete}
              totalLabel={t("groups.totalLabel")}
              createdAtLabel={t("groups.createdAt")}
              manageLabel={t("groups.manageMembers")}
              deleteLabel={t("groups.deleteAction")}
            />
          )}
        </Stack>
      </PageCard>

      <Dialog
        open={createDialogOpen}
        onClose={() => setCreateDialogOpen(false)}
        fullWidth
        maxWidth="sm"
      >
        <DialogTitle>{t("groups.createDialogTitle")}</DialogTitle>
        <DialogContent>
          <TextField
            label={t("groups.createDialogNameLabel")}
            value={groupName}
            onChange={(event) => {
              setGroupName(event.target.value);
              if (createError) {
                setCreateError("");
              }
            }}
            error={Boolean(createError)}
            helperText={createError}
            fullWidth
            required
            sx={{ mt: 1 }}
          />
        </DialogContent>
        <DialogActions>
          <Button
            onClick={() => setCreateDialogOpen(false)}
            color="inherit"
            disabled={creating}
          >
            {t("common.cancel")}
          </Button>
          <Button
            variant="contained"
            onClick={handleCreateGroup}
            disabled={creating}
          >
            {creating
              ? t("groups.createDialogSubmitting")
              : t("groups.createDialogSubmit")}
          </Button>
        </DialogActions>
      </Dialog>

      <ConfirmDialog
        open={Boolean(groupToDelete)}
        title={t("groups.deleteConfirmTitle")}
        description={
          <Stack spacing={1}>
            <Typography>{t("groups.deleteConfirmMessage")}</Typography>
            {groupToDelete && (
              <Typography fontWeight={600}>{groupToDelete.name}</Typography>
            )}
          </Stack>
        }
        confirmLabel={t("groups.deleteConfirmAction")}
        cancelLabel={t("common.cancel")}
        onConfirm={confirmDeleteGroup}
        onClose={() => setGroupToDelete(null)}
        confirming={deleting}
      />
    </CenteredContent>
  );
};

export default Groups;
