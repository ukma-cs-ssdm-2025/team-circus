import { useEffect, useMemo, useState } from "react";
import {
	Alert,
	Box,
	Button,
	Dialog,
	DialogActions,
	DialogContent,
	DialogTitle,
	FormControl,
	InputLabel,
	MenuItem,
	Select,
	Stack,
	TextField,
	Typography,
} from "@mui/material";
import { useLocation } from "react-router-dom";
import {
	CenteredContent,
	PageCard,
	PageHeader,
	ErrorAlert,
	LoadingSpinner,
	DocumentFilters,
	DocumentsGrid,
	ConfirmDialog,
} from "../components";
import { useLanguage } from "../contexts/LanguageContext";
import { useApi } from "../hooks";
import { API_ENDPOINTS } from "../constants";
import { createDocument, deleteDocument } from "../services";
import type {
	BaseComponentProps,
	DocumentsResponse,
	DocumentItem,
	GroupsResponse,
	DocumentFilters as DocumentFiltersType,
	GroupOption,
} from "../types";

type DocumentsProps = BaseComponentProps;

const Documents = ({ className = "" }: DocumentsProps) => {
	const { t } = useLanguage();
	const location = useLocation();
	const [filters, setFilters] = useState<DocumentFiltersType>({
		selectedGroup: "all",
		searchTerm: "",
	});
	const [createDialogOpen, setCreateDialogOpen] = useState(false);
	const [createForm, setCreateForm] = useState({
		name: "",
		group_uuid: "",
		content: "",
	});
	const [createErrors, setCreateErrors] = useState<Record<string, string>>({});
	const [createLoading, setCreateLoading] = useState(false);
	const [feedback, setFeedback] = useState<{
		type: "success" | "error";
		message: string;
	} | null>(null);
	const [documentToDelete, setDocumentToDelete] = useState<DocumentItem | null>(
		null,
	);
	const [deleting, setDeleting] = useState(false);

	const {
		data: documentsData,
		loading: documentsLoading,
		error: documentsError,
		refetch: refetchDocuments,
	} = useApi<DocumentsResponse>(API_ENDPOINTS.DOCUMENTS.BASE);

	const { data: groupsData, loading: groupsLoading } = useApi<GroupsResponse>(
		API_ENDPOINTS.GROUPS.BASE,
	);

	const documents = useMemo(
		() => documentsData?.documents ?? [],
		[documentsData],
	);
	const groups = useMemo(() => groupsData?.groups ?? [], [groupsData]);

	const groupOptions = useMemo((): GroupOption[] => {
		return groups.map((group) => ({
			value: group.uuid,
			label: group.name,
		}));
	}, [groups]);

	const groupNameByUUID = useMemo(() => {
		return groups.reduce<Record<string, string>>((acc, group) => {
			acc[group.uuid] = group.name;
			return acc;
		}, {});
	}, [groups]);

	const filteredDocuments = useMemo(() => {
		return documents.filter((document) => {
			const matchesGroup =
				filters.selectedGroup === "all" ||
				document.group_uuid === filters.selectedGroup;
			const matchesSearch = document.name
				.toLowerCase()
				.includes(filters.searchTerm.toLowerCase());
			return matchesGroup && matchesSearch;
		});
	}, [documents, filters.selectedGroup, filters.searchTerm]);

	const handleGroupChange = (value: string) => {
		setFilters((prev) => ({ ...prev, selectedGroup: value }));
	};

	const handleSearchChange = (value: string) => {
		setFilters((prev) => ({ ...prev, searchTerm: value }));
	};

	const isLoading = documentsLoading || groupsLoading;
	const canCreateDocument = groups.length > 0;

	useEffect(() => {
		const params = new URLSearchParams(location.search);
		const groupFromQuery = params.get("group");

		if (!groupFromQuery) {
			return;
		}

		const isKnownGroup = groups.some((group) => group.uuid === groupFromQuery);
		if (isKnownGroup) {
			setFilters((prev) => ({ ...prev, selectedGroup: groupFromQuery }));
		}
	}, [location.search, groups]);

	const validateCreateForm = () => {
		const errors: Record<string, string> = {};
		if (!createForm.name.trim()) {
			errors.name = t("documents.fieldRequired");
		}
		if (!createForm.group_uuid) {
			errors.group_uuid = t("documents.fieldRequired");
		}
		if (!createForm.content.trim()) {
			errors.content = t("documents.fieldRequired");
		}
		setCreateErrors(errors);
		return Object.keys(errors).length === 0;
	};

	const handleCreateDocument = async () => {
		if (!validateCreateForm()) {
			return;
		}

		setCreateLoading(true);
		try {
			await createDocument({
				...createForm,
				name: createForm.name.trim(),
				content: createForm.content.trim(),
			});
			setFeedback({ type: "success", message: t("documents.createSuccess") });
			setCreateDialogOpen(false);
			setCreateForm({ name: "", group_uuid: "", content: "" });
			setCreateErrors({});
			await refetchDocuments();
		} catch (error) {
			console.error("Failed to create document", error);
			setFeedback({ type: "error", message: t("documents.createError") });
		} finally {
			setCreateLoading(false);
		}
	};

	const handleDeleteDocument = (document: DocumentItem) => {
		setDocumentToDelete(document);
	};

	const confirmDeleteDocument = async () => {
		if (!documentToDelete) {
			return;
		}

		setDeleting(true);
		try {
			await deleteDocument(documentToDelete.uuid);
			setFeedback({ type: "success", message: t("documents.deleteSuccess") });
			setDocumentToDelete(null);
			await refetchDocuments();
		} catch (error) {
			console.error("Failed to delete document", error);
			setFeedback({ type: "error", message: t("documents.deleteError") });
		} finally {
			setDeleting(false);
		}
	};

	const handleDialogClose = () => {
		if (createLoading) {
			return;
		}
		setCreateDialogOpen(false);
		setCreateErrors({});
	};

	return (
		<CenteredContent className={className}>
			<PageCard>
				<PageHeader
					title={t("documents.title")}
					subtitle={t("documents.subtitle")}
				/>

				<Stack spacing={3}>
					<Stack spacing={2}>
						<DocumentFilters
							filters={filters}
							groupOptions={groupOptions}
							onGroupChange={handleGroupChange}
							onSearchChange={handleSearchChange}
							filterGroupLabel={t("documents.filterGroup")}
							filterAllLabel={t("documents.filterAll")}
							searchPlaceholder={t("documents.searchPlaceholder")}
						/>

						<Stack
							direction={{ xs: "column", sm: "row" }}
							justifyContent="flex-end"
						>
							<Button
								variant="contained"
								onClick={() => setCreateDialogOpen(true)}
								disabled={!canCreateDocument}
							>
								{t("documents.createButton")}
							</Button>
						</Stack>
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

					{isLoading && <LoadingSpinner py={6} />}

					{documentsError && (
						<ErrorAlert
							message={t("documents.error")}
							onRetry={refetchDocuments}
							retryText={t("documents.refresh")}
						/>
					)}

					{!isLoading && !documentsError && filteredDocuments.length === 0 && (
						<Typography color="text.secondary" align="center">
							{t("documents.empty")}
						</Typography>
					)}

					<DocumentsGrid
						documents={filteredDocuments}
						groupNameByUUID={groupNameByUUID}
						createdAtLabel={t("documents.createdAt")}
						noContentLabel={t("documents.noContent")}
						groupUnknownLabel={t("documents.groupUnknown")}
						deleteLabel={t("documents.deleteAction")}
						onDeleteDocument={handleDeleteDocument}
					/>
				</Stack>
			</PageCard>

			<Dialog
				open={createDialogOpen}
				onClose={handleDialogClose}
				fullWidth
				maxWidth="sm"
			>
				<DialogTitle>{t("documents.createDialogTitle")}</DialogTitle>
				<DialogContent>
					{!canCreateDocument && (
						<Alert severity="info" sx={{ mb: 2 }}>
							{t("documents.noGroupsHelper")}
						</Alert>
					)}
					<Stack spacing={2} py={1}>
						<TextField
							label={t("documents.createDialogNameLabel")}
							value={createForm.name}
							onChange={(event) =>
								setCreateForm((prev) => ({ ...prev, name: event.target.value }))
							}
							error={Boolean(createErrors.name)}
							helperText={createErrors.name}
							required
						/>
						<FormControl fullWidth error={Boolean(createErrors.group_uuid)}>
							<InputLabel>{t("documents.createDialogGroupLabel")}</InputLabel>
							<Select
								value={createForm.group_uuid}
								label={t("documents.createDialogGroupLabel")}
								onChange={(event) =>
									setCreateForm((prev) => ({
										...prev,
										group_uuid: event.target.value,
									}))
								}
							>
								{groupOptions.map((option) => (
									<MenuItem key={option.value} value={option.value}>
										{option.label}
									</MenuItem>
								))}
							</Select>
							{createErrors.group_uuid && (
								<Typography variant="caption" color="error" sx={{ mt: 0.5 }}>
									{createErrors.group_uuid}
								</Typography>
							)}
						</FormControl>

						<TextField
							label={t("documents.createDialogContentLabel")}
							value={createForm.content}
							onChange={(event) =>
								setCreateForm((prev) => ({
									...prev,
									content: event.target.value,
								}))
							}
							error={Boolean(createErrors.content)}
							helperText={createErrors.content}
							multiline
							minRows={4}
							required
						/>
					</Stack>
				</DialogContent>
				<DialogActions>
					<Button
						onClick={handleDialogClose}
						color="inherit"
						disabled={createLoading}
					>
						{t("common.cancel")}
					</Button>
					<Button
						variant="contained"
						onClick={handleCreateDocument}
						disabled={createLoading || !canCreateDocument}
					>
						{createLoading
							? t("documents.createDialogSubmitting")
							: t("documents.createDialogSubmit")}
					</Button>
				</DialogActions>
			</Dialog>

			<ConfirmDialog
				open={Boolean(documentToDelete)}
				title={t("documents.deleteConfirmTitle")}
				description={
					<Box>
						{t("documents.deleteConfirmMessage")}
						{documentToDelete && (
							<Typography
								component="span"
								fontWeight={600}
								display="block"
								mt={1}
							>
								{documentToDelete.name}
							</Typography>
						)}
					</Box>
				}
				confirmLabel={t("documents.deleteConfirmAction")}
				cancelLabel={t("common.cancel")}
				onConfirm={confirmDeleteDocument}
				onClose={() => setDocumentToDelete(null)}
				confirming={deleting}
			/>
		</CenteredContent>
	);
};

export default Documents;
