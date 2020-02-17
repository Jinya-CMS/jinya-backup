using Hangfire.Dashboard;

namespace Jinya.Backup.Auth
{
    public class JinyaDashboardAuthorization : IDashboardAuthorizationFilter
    {
        public bool Authorize(DashboardContext context)
        {
            return true;
        }
    }
}