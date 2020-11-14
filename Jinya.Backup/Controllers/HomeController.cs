using Microsoft.AspNetCore.Mvc;

namespace Jinya.Backup.Controllers
{
    public class HomeController : Controller
    {
        public IActionResult Index()
        {
            return View();
        }
    }
}