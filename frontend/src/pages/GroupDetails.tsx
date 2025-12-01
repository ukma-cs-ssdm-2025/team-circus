import {
	Alert,
	Autocomplete,
	Button,
	FormControl,
	InputLabel,
	MenuItem,
	Paper,
	Select,
	Stack,
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableRow,
	TextField,
	Typography,
} from "@mui/material";
import { useEffect, useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import {
	CenteredContent,
	ConfirmDialog,
	ErrorAlert,
	LoadingSpinner,
	PageCard,
	PageHeader,
} from "../components";
import { API_ENDPOINTS, MEMBER_ROLES, ROUTES } from "../constants";
import { useAuth } from "../contexts/AuthContextBase";
import { useLanguage } from "../contexts/LanguageContext";
import { useApi } from "../hooks";
import {
	addGroupMember,
	removeGroupMember,
	updateGroupMemberRole,
} from "../services";
import type {
	BaseComponentProps,
	GroupItem,
	MemberItem,
	MembersResponse,
	UsersResponse,
} from "../types";
import { formatDate } from "../utils";

type GroupDetailsProps = BaseComponentProps;

const roleOptions = MEMBER_ROLES.map((role) => ({
	value: role,
	key: `members.role.${role}`,
}));

const memberFormRoleOptions = roleOptions.filter(
	(option) => option.value !== MEMBER_ROLES[0],
);

const GroupDetails = ({ className = "" }: GroupDetailsProps) => {
	const { t } = useLanguage();
	const { user } = useAuth();
	const navigate = useNavigate();
	const { uuid } = useParams<{ uuid: string }>();
	const groupUUID = uuid ?? "";

	const [memberForm, setMemberForm] = useState({
		user_uuid: "",
		role: MEMBER_ROLES[2],
	});
	const [memberFormErrors, setMemberFormErrors] = useState<
		Record<string, string>
	>({});
	const [memberFeedback, setMemberFeedback] = useState<{
		type: "success" | "error";
		message: string;
	} | null>(null);
	const [addingMember, setAddingMember] = useState(false);
	const [updatingMemberId, setUpdatingMemberId] = useState<string | null>(null);
	const [removingMemberId, setRemovingMemberId] = useState<string | null>(null);
	const [roleConfirmation, setRoleConfirmation] = useState<{
		memberUUID: string;
		role: MemberRole;
	} | null>(null);

	const {
		data: groupData,
		loading: groupLoading,
		error: groupError,
		refetch: refetchGroup,
	} = useApi<GroupItem>(API_ENDPOINTS.GROUPS.DETAIL(groupUUID));

	const {
		data: membersData,
		loading: membersLoading,
		error: membersError,
		refetch: refetchMembers,
	} = useApi<MembersResponse>(API_ENDPOINTS.GROUPS.MEMBERS(groupUUID));

	const { data: usersData, loading: usersLoading } = useApi<UsersResponse>(
		API_ENDPOINTS.USERS.BASE,
	);

	const members = useMemo(() => membersData?.members ?? [], [membersData]);
	const users = useMemo(() => usersData?.users ?? [], [usersData]);

	const selectableUsers = useMemo(
		() =>
			membersLoading
				? []
				: users.filter(
						(candidate) =>
							candidate.uuid !== user?.uuid &&
							!members.some((member) => member.user_uuid === candidate.uuid),
					),
		[members, membersLoading, user?.uuid, users],
	);

	const userListLoading = usersLoading || membersLoading;

	const selectableUserOptions = useMemo(
		() =>
			selectableUsers.map((userOption) => ({
				value: userOption.uuid,
				label: `${userOption.login} (${userOption.email})`,
			})),
		[selectableUsers],
	);

	const userByUUID = useMemo(() => {
		return users.reduce<Record<string, { login: string; email: string }>>(
			(acc, current) => {
				acc[current.uuid] = { login: current.login, email: current.email };
				return acc;
			},
			{},
		);
	}, [users]);

	const currentMember = useMemo(() => {
		if (!user) {
			return undefined;
		}

		return (
			members.find((member) => member.user_uuid === user.uuid) ??
			members.find(
				(member) => userByUUID[member.user_uuid]?.login === user.login,
			)
		);
	}, [members, user, userByUUID]);

	const canManageMembers = Boolean(currentMember?.role === MEMBER_ROLES[0]);

	const validateMemberForm = () => {
		const errors: Record<string, string> = {};
		if (!memberForm.user_uuid) {
			errors.user_uuid = t("groupDetails.fieldRequired");
		}
		if (!memberForm.role) {
			errors.role = t("groupDetails.fieldRequired");
		}
		setMemberFormErrors(errors);
		return Object.keys(errors).length === 0;
	};

	useEffect(() => {
		if (memberForm.user_uuid && memberForm.user_uuid === user?.uuid) {
			setMemberForm((prev) => ({ ...prev, user_uuid: "" }));
		}
	}, [memberForm.user_uuid, user?.uuid]);

	const handleAddMember = async () => {
		if (!groupUUID || !validateMemberForm()) {
			return;
		}

		setAddingMember(true);
		try {
			await addGroupMember(groupUUID, memberForm);
			setMemberFeedback({
				type: "success",
				message: t("groupDetails.addSuccess"),
			});
			setMemberForm({ user_uuid: "", role: MEMBER_ROLES[2] });
			setMemberFormErrors({});
			await refetchMembers();
		} catch (error) {
			console.error("Failed to add member", error);
			setMemberFeedback({ type: "error", message: t("groupDetails.addError") });
		} finally {
			setAddingMember(false);
		}
	};

	const handleRoleChange = async (memberUUID: string, role: MemberRole) => {
		if (!groupUUID) {
			return;
		}

		setUpdatingMemberId(memberUUID);
		try {
			await updateGroupMemberRole(groupUUID, memberUUID, {
				role: role as MemberItem["role"],
			});
			setMemberFeedback({
				type: "success",
				message: t("groupDetails.updateSuccess"),
			});
			await refetchMembers();
		} catch (error) {
			console.error("Failed to change member role", error);
			setMemberFeedback({
				type: "error",
				message: t("groupDetails.updateError"),
			});
		} finally {
			setUpdatingMemberId(null);
		}
	};

	const handleConfirmRoleChange = async () => {
		if (!roleConfirmation) {
			return;
		}

		await handleRoleChange(roleConfirmation.memberUUID, roleConfirmation.role);
		setRoleConfirmation(null);
	};

	const handleRemoveMember = async (memberUUID: string) => {
		if (!groupUUID) {
			return;
		}

		setRemovingMemberId(memberUUID);
		try {
			await removeGroupMember(groupUUID, memberUUID);
			setMemberFeedback({
				type: "success",
				message: t("groupDetails.removeSuccess"),
			});
			await refetchMembers();
		} catch (error) {
			console.error("Failed to remove member", error);
			setMemberFeedback({
				type: "error",
				message: t("groupDetails.removeError"),
			});
		} finally {
			setRemovingMemberId(null);
		}
	};

	return (
		<CenteredContent className={className}>
			<PageCard>
				<Stack spacing={3}>
					<Stack
						direction={{ xs: "column", sm: "row" }}
						justifyContent="space-between"
						gap={2}
					>
						<PageHeader
							title={groupData?.name || t("groupDetails.titleFallback")}
							subtitle={t("groupDetails.subtitle")}
						/>
						<Button variant="outlined" onClick={() => navigate(ROUTES.GROUPS)}>
							{t("groupDetails.backToGroups")}
						</Button>
					</Stack>

					{(groupLoading || membersLoading) && <LoadingSpinner />}

					{groupError && (
						<ErrorAlert
							message={t("groupDetails.loadError")}
							onRetry={refetchGroup}
							retryText={t("groupDetails.retry")}
						/>
					)}

					{membersError && (
						<ErrorAlert
							message={t("groupDetails.membersError")}
							onRetry={refetchMembers}
							retryText={t("groupDetails.retry")}
						/>
					)}

					{groupData && (
						<Paper
							elevation={0}
							sx={{
								p: 3,
								borderRadius: 3,
								border: "1px solid",
								borderColor: "divider",
							}}
						>
							<Typography variant="h6" fontWeight={700}>
								{groupData.name}
							</Typography>
							<Typography variant="body2" color="text.secondary">
								{`${t("groupDetails.createdAtLabel")}: ${formatDate(groupData.created_at)}`}
							</Typography>
						</Paper>
					)}

					{memberFeedback && (
						<Alert
							severity={memberFeedback.type}
							onClose={() => setMemberFeedback(null)}
							variant="outlined"
						>
							{memberFeedback.message}
						</Alert>
					)}

					<Paper
						elevation={0}
						sx={{
							p: 3,
							borderRadius: 3,
							border: "1px solid",
							borderColor: "divider",
						}}
					>
						<Typography variant="h6" fontWeight={700} sx={{ mb: 2 }}>
							{t("groupDetails.addMemberTitle")}
						</Typography>
						{!canManageMembers && (
							<Alert severity="info" sx={{ mb: 2 }}>
								{t("groupDetails.permissionsHint")}
							</Alert>
						)}
						<Stack spacing={2}>
							<Autocomplete
								options={selectableUserOptions}
								loading={userListLoading}
								isOptionEqualToValue={(option, value) =>
									option.value === value.value
								}
								value={
									memberForm.user_uuid
										? {
												value: memberForm.user_uuid,
												label:
													userByUUID[memberForm.user_uuid]?.login ||
													memberForm.user_uuid,
											}
										: null
								}
								onChange={(_, option) => {
									setMemberForm((prev) => ({
										...prev,
										user_uuid: option?.value ?? "",
									}));
									if (memberFormErrors.user_uuid) {
										setMemberFormErrors((prev) => ({ ...prev, user_uuid: "" }));
									}
								}}
								renderInput={(params) => (
									<TextField
										{...params}
										label={t("groupDetails.userFieldLabel")}
										helperText={memberFormErrors.user_uuid}
										error={Boolean(memberFormErrors.user_uuid)}
									/>
								)}
								disabled={!canManageMembers || userListLoading}
							/>

							<FormControl fullWidth error={Boolean(memberFormErrors.role)}>
								<InputLabel>{t("groupDetails.roleFieldLabel")}</InputLabel>
								<Select
									value={memberForm.role}
									label={t("groupDetails.roleFieldLabel")}
									onChange={(event) => {
										setMemberForm((prev) => ({
											...prev,
											role: event.target.value as MemberItem["role"],
										}));
										if (memberFormErrors.role) {
											setMemberFormErrors((prev) => ({ ...prev, role: "" }));
										}
									}}
									disabled={!canManageMembers}
								>
									{memberFormRoleOptions.map((option) => (
										<MenuItem key={option.value} value={option.value}>
											{t(option.key)}
										</MenuItem>
									))}
								</Select>
								{memberFormErrors.role && (
									<Typography variant="caption" color="error" sx={{ mt: 0.5 }}>
										{memberFormErrors.role}
									</Typography>
								)}
							</FormControl>

							<Stack direction="row" justifyContent="flex-end">
								<Button
									variant="contained"
									onClick={handleAddMember}
									disabled={!canManageMembers || addingMember}
								>
									{addingMember
										? t("groupDetails.addMemberSubmitting")
										: t("groupDetails.addMemberSubmit")}
								</Button>
							</Stack>
						</Stack>
					</Paper>

					<Paper
						elevation={0}
						sx={{
							p: 3,
							borderRadius: 3,
							border: "1px solid",
							borderColor: "divider",
						}}
					>
						<Typography variant="h6" fontWeight={700} sx={{ mb: 2 }}>
							{t("groupDetails.membersTitle")}
						</Typography>
						{members.length === 0 ? (
							<Typography color="text.secondary">
								{t("groupDetails.membersEmpty")}
							</Typography>
						) : (
							<Table size="small">
								<TableHead>
									<TableRow>
										<TableCell>{t("groupDetails.memberColumnUser")}</TableCell>
										<TableCell>{t("groupDetails.memberColumnRole")}</TableCell>
										<TableCell align="right">
											{t("groupDetails.memberColumnActions")}
										</TableCell>
									</TableRow>
								</TableHead>
								<TableBody>
									{members.map((member) => (
										<TableRow key={member.user_uuid}>
											<TableCell>
												<Typography fontWeight={600}>
													{userByUUID[member.user_uuid]?.login ||
														member.user_uuid}
												</Typography>
												<Typography variant="caption" color="text.secondary">
													{member.user_uuid}
												</Typography>
											</TableCell>
											<TableCell sx={{ minWidth: 180 }}>
												{member.role === MEMBER_ROLES[0] ? (
													<Typography fontWeight={600}>
														{t(`members.role.${member.role}`)}
													</Typography>
												) : (
													<FormControl
														size="small"
														fullWidth
														disabled={
															!canManageMembers ||
															updatingMemberId === member.user_uuid
														}
													>
														<InputLabel>
															{t("groupDetails.roleFieldLabel")}
														</InputLabel>
														<Select
															value={member.role}
															label={t("groupDetails.roleFieldLabel")}
															onChange={(event) => {
																const nextRole = event.target
																	.value as MemberRole;
																if (
																	nextRole === MEMBER_ROLES[0] &&
																	member.role !== MEMBER_ROLES[0]
																) {
																	setRoleConfirmation({
																		memberUUID: member.user_uuid,
																		role: nextRole,
																	});
																	return;
																}
																void handleRoleChange(
																	member.user_uuid,
																	nextRole,
																);
															}}
														>
															{roleOptions.map((option) => (
																<MenuItem
																	key={option.value}
																	value={option.value}
																>
																	{t(option.key)}
																</MenuItem>
															))}
														</Select>
													</FormControl>
												)}
											</TableCell>
											<TableCell align="right">
												<Button
													color="error"
													size="small"
													onClick={() => handleRemoveMember(member.user_uuid)}
													disabled={
														!canManageMembers ||
														member.role === MEMBER_ROLES[0] ||
														removingMemberId === member.user_uuid
													}
												>
													{removingMemberId === member.user_uuid
														? t("groupDetails.removeMemberSubmitting")
														: t("groupDetails.removeMemberAction")}
												</Button>
											</TableCell>
										</TableRow>
									))}
								</TableBody>
							</Table>
						)}
					</Paper>
				</Stack>
			</PageCard>
			<ConfirmDialog
				open={Boolean(roleConfirmation)}
				title={t("groupDetails.promoteConfirmTitle")}
				description={t("groupDetails.promoteConfirmMessage")}
				confirmLabel={t("groupDetails.promoteConfirmAction")}
				cancelLabel={t("common.cancel")}
				onConfirm={handleConfirmRoleChange}
				onClose={() => setRoleConfirmation(null)}
			/>
		</CenteredContent>
	);
};

export default GroupDetails;
