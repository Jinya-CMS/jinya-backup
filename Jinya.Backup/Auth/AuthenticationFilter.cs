using System.Linq;
using System.Threading.Tasks;
using Jinya.Backup.Controllers;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.Filters;

namespace Jinya.Backup.Auth
{
    public class AuthenticationFilter : IAsyncActionFilter
    {
        public async Task OnActionExecutionAsync(ActionExecutingContext context, ActionExecutionDelegate next)
        {
            if (context.Controller is HomeController || context.HttpContext.Request.Headers["AuthKey"].First() == AuthenticationKeyStorage.AuthKey)
            {
                await next();
            }
            else
            {
                context.Result = new OkResult();
            }
        }
    }
}