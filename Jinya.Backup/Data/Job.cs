using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;

namespace Jinya.Backup.Data
{
    public class Job
    {
        public Job()
        {
            Steps = new List<Step>();
        }

        [Key] public Guid Id { get; set; } = Guid.NewGuid();
        public string Name { get; set; }

        public ICollection<Step> Steps { get; set; }
    }
}