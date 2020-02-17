using Microsoft.AspNetCore.Mvc;

namespace Jinya.Backup.Controllers
{
    public class HomeController : Controller
    {
        // GET
        public IActionResult Index()
        {
            return View();
        }
    }
}