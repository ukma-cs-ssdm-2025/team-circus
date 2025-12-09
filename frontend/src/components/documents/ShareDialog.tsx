import ContentCopyIcon from "@mui/icons-material/ContentCopy";
import CheckIcon from "@mui/icons-material/Check";
import {
	Alert,
	Box,
	Button,
	Chip,
	CircularProgress,
	Dialog,
	DialogActions,
	DialogContent,
	DialogTitle,
	FormControl,
	IconButton,
	InputAdornment,
	InputLabel,
	MenuItem,
	Select,
	Stack,
	TextField,
	Typography,
} from "@mui/material";
import { useEffect, useMemo, useState } from "react";
import { useLanguage } from "../../contexts/LanguageContext";
import type { ShareLinkResponse } from "../../types";

type ShareDialogProps = {
	open: boolean;
	loading: boolean;
	link: ShareLinkResponse | null;
	error?: string | null;
	defaultExpirationDays: number;
	onClose: () => void;
	onGenerate: (expirationDays: number) => Promise<void> | void;
};

const EXPIRATION_OPTIONS = [1, 7, 14, 30, 60, 90];

export const ShareDialog = ({
	open,
	loading,
	link,
	error,
	defaultExpirationDays,
	onClose,
	onGenerate,
}: ShareDialogProps) => {
	const { t } = useLanguage();
	const [expiration, setExpiration] = useState(defaultExpirationDays);
	const [copied, setCopied] = useState(false);

	useEffect(() => {
		if (open) {
			setExpiration(defaultExpirationDays);
			setCopied(false);
		}
	}, [defaultExpirationDays, open]);

	const expiresLabel = useMemo(() => {
		if (!link?.expires_at) {
			return "";
		}
		return new Intl.DateTimeFormat(undefined, {
			year: "numeric",
			month: "short",
			day: "numeric",
			hour: "2-digit",
			minute: "2-digit",
		}).format(new Date(link.expires_at));
	}, [link?.expires_at]);

	const handleCopy = async () => {
		if (!link?.url) return;
		try {
			await navigator.clipboard.writeText(link.url);
			setCopied(true);
			setTimeout(() => setCopied(false), 1600);
		} catch (copyError) {
			console.error("Failed to copy link", copyError);
			setCopied(false);
		}
	};

	const handleGenerate = async () => {
		setCopied(false);
		await onGenerate(expiration);
	};

	return (
		<Dialog open={open} onClose={onClose} fullWidth maxWidth="sm">
			<DialogTitle>
				{t("shareDialog.title")}
				<Typography variant="body2" color="text.secondary" mt={1}>
					{t("shareDialog.subtitle")}
				</Typography>
			</DialogTitle>

			<DialogContent
				sx={{
					background:
						"linear-gradient(135deg, rgba(37,99,235,0.05), rgba(14,165,233,0.08))",
				}}
			>
				<Stack spacing={2} mt={1}>
					<FormControl fullWidth>
						<InputLabel>{t("shareDialog.expirationLabel")}</InputLabel>
						<Select
							value={expiration}
							label={t("shareDialog.expirationLabel")}
							onChange={(event) =>
								setExpiration(Number(event.target.value) || defaultExpirationDays)
							}
						>
							{EXPIRATION_OPTIONS.map((option) => (
								<MenuItem key={option} value={option}>
									{`${option} ${t("shareDialog.daysLabel")}`}
								</MenuItem>
							))}
						</Select>
						{/* Removed expirationHint text */}
					</FormControl>

					{error && (
						<Alert severity="error" variant="outlined">
							{error}
						</Alert>
					)}

					{link && (
						<Stack spacing={1.5}>
							<TextField
								label={t("shareDialog.linkLabel")}
								value={link.url}
								InputProps={{
									readOnly: true,
									endAdornment: (
										<InputAdornment position="end">
											<IconButton
												onClick={handleCopy}
												edge="end"
												aria-label="copy link"
											>
												{copied ? (
													<CheckIcon color="success" fontSize="small" />
												) : (
													<ContentCopyIcon fontSize="small" />
												)}
											</IconButton>
										</InputAdornment>
									),
								}}
							/>
							<Box display="flex" gap={1} alignItems="center" flexWrap="wrap">
								<Chip
									color="primary"
									label={t("shareDialog.readOnlyTag")}
									variant="filled"
									size="small"
									sx={{ fontWeight: 700 }}
								/>
								{expiresLabel && (
									<Typography variant="body2" color="text.secondary">
										{`${t("shareDialog.expiresAtLabel")}: ${expiresLabel}`}
									</Typography>
								)}
							</Box>
						</Stack>
					)}
				</Stack>
			</DialogContent>

			<DialogActions sx={{ px: 3, pb: 2, pt: 1 }}>
				<Button onClick={onClose} color="inherit">
					{t("common.cancel")}
				</Button>
				<Button
					variant="contained"
					onClick={handleGenerate}
					disabled={loading}
					startIcon={
						loading ? <CircularProgress size={16} color="inherit" /> : null
					}
				>
					{loading
						? t("shareDialog.generating")
						: t("shareDialog.generateButton")}
				</Button>
			</DialogActions>
		</Dialog>
	);
};
