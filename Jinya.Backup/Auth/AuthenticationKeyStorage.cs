using System;
using System.Linq;
using System.Security.Cryptography;
using Hangfire;

namespace Jinya.Backup.Auth
{
    public class AuthenticationKeyStorage
    {
        public static string AuthKey { get; private set; }

        static AuthenticationKeyStorage()
        {
            CreateAuthKey();
        }

        public static void CreateAuthKey()
        {
            var sha512 = SHA512.Create();
            var authkey = Convert.ToBase64String(sha512.ComputeHash(Guid.NewGuid().ToByteArray())
                .Concat(sha512.ComputeHash(Guid.NewGuid().ToByteArray())).ToArray());
            AuthKey = authkey;
        }

        public static void Run()
        {
            Console.WriteLine($"Authentication key is {AuthKey}");
            var cron = Cron.Daily();
            var recurringJobManager = new RecurringJobManager(JobStorage.Current);
            recurringJobManager.AddOrUpdate("RecreateAuthKey", () => CreateAuthKey(), cron);
        }
    }
}