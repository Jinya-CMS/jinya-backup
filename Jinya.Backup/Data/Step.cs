using System;
using System.ComponentModel.DataAnnotations;

namespace Jinya.Backup.Data
{
    public enum StepType
    {
        Ssh = 1,
        FileTransfer = 2,
        Cleanup = 3
    }

    public class Step
    {
        [Key] public Guid Id { get; set; } = Guid.NewGuid();
        public int Position { get; set; }
        public StepType StepType { get; set; }
        public string Command { get; set; }
        public string Server { get; set; }
        public string Username { get; set; }
        public string Password { get; set; }
        public string Source { get; set; }
        public string Target { get; set; }

        public Job Job { get; set; }
    }
}