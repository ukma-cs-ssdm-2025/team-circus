import {
	Alert,
	Box,
	Chip,
	Container,
	Paper,
	Stack,
	Typography,
} from "@mui/material";
import { useEffect, useMemo, useState } from "react";
import { useLocation } from "react-router-dom";
import { LoadingSpinner, MarkdownPreview } from "../components";
import { useLanguage } from "../contexts/LanguageContext";
import { useCollaborativeEditor } from "../hooks";
import { fetchPublicDocument } from "../services";
import type { DocumentItem } from "../types";
import { useTheme } from "@mui/material/styles";

const formatDateTime = (value: string | number) => {
	const date =
		typeof value === "number"
			? new Date(value * 1000)
			: new Date(Date.parse(value));
	if (Number.isNaN(date.getTime())) {
		return "";
	}

	return new Intl.DateTimeFormat(undefined, {
		year: "numeric",
		month: "short",
		day: "numeric",
		hour: "2-digit",
		minute: "2-digit",
	}).format(date);
};

const PublicDocument = () => {
	const { t } = useLanguage();
	const theme = useTheme();
	const location = useLocation();
	const [document, setDocument] = useState<DocumentItem | null>(null);
	const [status, setStatus] = useState<number | null>(null);
	const [error, setError] = useState<string | null>(null);
	const [loading, setLoading] = useState(true);

	const params = useMemo(
		() => new URLSearchParams(location.search),
		[location.search],
	);

	const docParam = params.get("doc") ?? "";
	const sigParam = params.get("sig") ?? "";
	const expParam = params.get("exp") ?? "";
	const hasValidParams = Boolean(docParam && sigParam && expParam);
	const roomParams = useMemo(
		() => ({
			sig: sigParam,
			exp: expParam,
		}),
		[sigParam, expParam],
	);
	const guestUser = useMemo(
		() => ({
			id: `guest-${sigParam.slice(0, 8) || Math.random().toString(16).slice(2, 10)}`,
			name: "Guest",
			role: "viewer",
		}),
		[sigParam],
	);

	const { content: liveContent } = useCollaborativeEditor({
		documentId: docParam || "unknown",
		user: guestUser,
		wsPath: "/ws/public/documents",
		roomParams,
		enabled: hasValidParams,
	});
	const displayContent =
		(liveContent ?? "") !== "" ? liveContent : document?.content ?? "";

	const expiresLabel = useMemo(() => {
		if (!expParam) {
			return "";
		}
		return formatDateTime(Number(expParam));
	}, [expParam]);

	useEffect(() => {
		setLoading(true);
		setError(null);
		setStatus(null);
		setDocument(null);

		if (!docParam || !sigParam || !expParam) {
			setError(t("publicDocument.missingParams"));
			setStatus(400);
			setLoading(false);
			return;
		}

		fetchPublicDocument({
			doc: docParam,
			sig: sigParam,
			exp: expParam,
		})
			.then((data) => {
				setDocument(data);
			})
			.catch((err) => {
				const knownStatus =
					err && typeof err === "object" && "status" in err
						? Number((err as { status?: number }).status)
						: null;
				setStatus(Number.isFinite(knownStatus) ? knownStatus : null);
				setError(
					err instanceof Error ? err.message : t("publicDocument.error"),
				);
			})
			.finally(() => setLoading(false));
	}, [docParam, expParam, sigParam, t]);

	const createdLabel = document ? formatDateTime(document.created_at) : "";
	const isExpired = status === 410;
	const isInvalid = status === 404 || status === 400;

	const pageBackground =
		theme.palette.mode === "light"
			? "linear-gradient(135deg, #eef2ff 0%, #f8fafc 35%, #ffffff 100%)"
			: "radial-gradient(circle at 15% 20%, rgba(59,130,246,0.12), transparent 26%), radial-gradient(circle at 82% 12%, rgba(14,165,233,0.1), transparent 22%), linear-gradient(180deg, #0b1226 0%, #0b1020 50%, #050914 100%)";

	const cardBackground =
		theme.palette.mode === "light"
			? "linear-gradient(135deg, rgba(255,255,255,0.98), rgba(248,250,252,0.94))"
			: "linear-gradient(135deg, rgba(24, 33, 53, 0.98), rgba(14, 20, 35, 0.96))";

	const surfaceBackground =
		theme.palette.mode === "light"
			? "rgba(255,255,255,0.9)"
			: "rgba(15,23,42,0.85)";

	const surfaceBorder =
		theme.palette.mode === "light"
			? "rgba(148,163,184,0.35)"
			: "rgba(148,163,184,0.25)";

	return (
		<Box
			sx={{
				minHeight: "100vh",
				background: pageBackground,
				py: { xs: 4, md: 6 },
			}}
		>
			<Container maxWidth="md">
				<Paper
					elevation={6}
					sx={{
						p: { xs: 3, md: 4 },
						borderRadius: 3,
						background: cardBackground,
						color: theme.palette.text.primary,
					}}
				>
					<Stack spacing={2.5}>
						<Box display="flex" alignItems="center" gap={1.5} flexWrap="wrap">
							<Chip
								color="primary"
								label={t("shareDialog.readOnlyTag")}
								sx={{ fontWeight: 700 }}
							/>
							<Typography variant="body2" color="text.secondary">
								{t("publicDocument.readOnlyNote")}
							</Typography>
						</Box>

						{expiresLabel && (
							<Typography variant="body2" color="text.secondary">
								{`${t("publicDocument.expiryLabel")}: ${expiresLabel}`}
							</Typography>
						)}

						{loading && <LoadingSpinner py={4} />}

						{!loading && error && (
							<Alert severity={isExpired ? "warning" : "error"} variant="filled">
								<Typography variant="h6" gutterBottom>
									{isExpired
										? t("publicDocument.expiredTitle")
										: isInvalid
											? t("publicDocument.invalidTitle")
											: t("publicDocument.error")}
								</Typography>
								<Typography variant="body2">
									{isExpired
										? t("publicDocument.expiredMessage")
										: isInvalid
											? t("publicDocument.invalidMessage")
											: error}
								</Typography>
							</Alert>
						)}

						{!loading && !error && document && (
							<Stack spacing={2.5}>
								<Box>
									<Typography variant="h4" fontWeight={800} gutterBottom>
										{document.name || t("documentEditor.fallbackTitle")}
									</Typography>
									<Stack direction="row" spacing={2} flexWrap="wrap">
										{createdLabel && (
											<Typography variant="body2" color="text.secondary">
												{`${t("publicDocument.metaCreated")}: ${createdLabel}`}
											</Typography>
										)}
										{document.uuid && (
											<Typography variant="body2" color="text.secondary">
												{`${t("publicDocument.metaId")}: ${document.uuid}`}
											</Typography>
										)}
									</Stack>
								</Box>

								<Paper
									variant="outlined"
									sx={{
										p: { xs: 2, md: 3 },
										borderRadius: 2.5,
										borderColor: surfaceBorder,
										background: surfaceBackground,
									}}
								>
									<MarkdownPreview
										content={displayContent}
										emptyText={t("documentEditor.previewEmpty")}
									/>
								</Paper>
							</Stack>
						)}
					</Stack>
				</Paper>
			</Container>
		</Box>
	);
};

export default PublicDocument;
