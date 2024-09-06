import { NavigationSidebar } from "@/components/navigation/NavigationSidebar";

const MainLayout = async ({
    children
}: {
    children: React.ReactNode, servers: any
}) => {
    return (
        <div className="h-full">
            <div className="max-md:hidden md:flex h-full w-[73px] z-30 flex-col fixed inset-y-0">
                <NavigationSidebar />
            </div>
            <main className="md:pl-[72px] h-full">
                {children}
            </main>
        </div>
    );
}

export default MainLayout;