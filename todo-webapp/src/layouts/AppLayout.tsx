import type { JSX } from "react";
import { Outlet } from "react-router-dom";
import { AppShell, Header, Footer } from "@wso2/oxygen-ui";

// The wireframe draws a brand-only navbar and no sidebar — this is a
// single-screen app with no sign-in and no navigation beyond that one
// screen, so the shell carries no Sidebar and no account/user-menu cluster.
export default function AppLayout(): JSX.Element {
  return (
    <AppShell>
      <AppShell.Navbar>
        <Header minimal>
          <Header.Brand>
            <Header.BrandTitle>Todo App</Header.BrandTitle>
          </Header.Brand>
          <Header.Spacer />
        </Header>
      </AppShell.Navbar>

      <AppShell.Main>
        <Outlet />
      </AppShell.Main>

      <AppShell.Footer>
        <Footer>
          <Footer.Copyright>© {new Date().getFullYear()} WSO2 LLC.</Footer.Copyright>
        </Footer>
      </AppShell.Footer>
    </AppShell>
  );
}
