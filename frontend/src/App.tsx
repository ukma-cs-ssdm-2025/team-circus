import {
	Outlet,
	Route,
	BrowserRouter as Router,
	Routes,
} from "react-router-dom";
import { Layout } from "./components/layout";
import { ProtectedRoute } from "./components/ProtectedRoute";
import { ROUTES } from "./constants";
import { AuthProvider } from "./contexts/AuthContext";
import { LanguageProvider } from "./contexts/LanguageContext";
import { ThemeProvider } from "./contexts/ThemeContext";
import {
	AuthTest,
	DocumentEditor,
	Documents,
	GroupDetails,
	Groups,
	Home,
	Login,
	NotFound,
	Register,
	Settings,
} from "./pages";

const stripLeadingSlash = (path: string) => path.replace(/^\//, "");

const ProtectedLayout = () => (
	<Layout>
		<Outlet />
	</Layout>
);

function App() {
	return (
		<ThemeProvider>
			<LanguageProvider>
				<AuthProvider>
					<Router>
						<Routes>
							{/* Public routes */}
							<Route path={ROUTES.LOGIN} element={<Login />} />
							<Route path={ROUTES.REGISTER} element={<Register />} />

							{/* Protected routes */}
							<Route element={<ProtectedRoute />}>
								<Route element={<ProtectedLayout />}>
									<Route index element={<Home />} />
									<Route
										path={stripLeadingSlash(ROUTES.DOCUMENTS)}
										element={<Documents />}
									/>
									<Route
										path={`${stripLeadingSlash(ROUTES.DOCUMENTS)}/:uuid`}
										element={<DocumentEditor />}
									/>
									<Route
										path={stripLeadingSlash(ROUTES.GROUPS)}
										element={<Groups />}
									/>
									<Route
										path={`${stripLeadingSlash(ROUTES.GROUPS)}/:uuid`}
										element={<GroupDetails />}
									/>
									<Route
										path={stripLeadingSlash(ROUTES.SETTINGS)}
										element={<Settings />}
									/>
									<Route path="auth-test" element={<AuthTest />} />
									<Route path="*" element={<NotFound />} />
								</Route>
							</Route>
						</Routes>
					</Router>
				</AuthProvider>
			</LanguageProvider>
		</ThemeProvider>
	);
}

export default App;
