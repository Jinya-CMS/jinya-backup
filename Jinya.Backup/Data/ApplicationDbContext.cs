using Microsoft.EntityFrameworkCore;

namespace Jinya.Backup.Data
{
    public class ApplicationDbContext : DbContext
    {
        public DbSet<Job> Jobs { get; set; }
        public DbSet<Step> Steps { get; set; }

        public ApplicationDbContext(DbContextOptions options) : base(options)
        {
        }
        
        protected override void OnModelCreating(ModelBuilder modelBuilder)
        {
            modelBuilder.Entity<Job>()
                .HasMany(j => j.Steps)
                .WithOne(s => s.Job)
                .IsRequired()
                .OnDelete(DeleteBehavior.Cascade);
        }
    }
}