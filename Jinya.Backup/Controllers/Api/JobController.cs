using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Threading.Tasks;
using Hangfire;
using Jinya.Backup.Data;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using Renci.SshNet;

namespace Jinya.Backup.Controllers.Api
{
    [ApiController]
    [Route("api/job")]
    public class JobController : ControllerBase
    {
        private readonly ApplicationDbContext _dbContext;
        private readonly IRecurringJobManager _recurringJobManager;

        public JobController(ApplicationDbContext dbContext, IRecurringJobManager recurringJobManager)
        {
            _dbContext = dbContext;
            _recurringJobManager = recurringJobManager;
        }

        [HttpGet]
        public async Task<IEnumerable<dynamic>> GetAll()
        {
            return await _dbContext.Jobs.AsNoTracking().Select(j => new {j.Id, j.Name}).ToListAsync();
        }

        [HttpGet("{id}")]
        public async Task<IActionResult> GetById(Guid id)
        {
            try
            {
                var dbJob = await _dbContext.Jobs.SingleOrDefaultAsync(j => j.Id == id);
                if (dbJob == null)
                {
                    return NotFound();
                }

                return Ok(dbJob);
            }
            catch (Exception ex)
            {
                return BadRequest(ex.Message);
            }
        }

        [HttpPost]
        public async Task<IActionResult> PostNew([FromBody] Job job)
        {
            try
            {
                await _dbContext.Jobs.AddAsync(job);
                await _dbContext.SaveChangesAsync();
                return NoContent();
            }
            catch (Exception ex)
            {
                return BadRequest(ex.Message);
            }
        }

        [HttpPut("{id}")]
        public async Task<IActionResult> Put([FromRoute] Guid id, [FromBody] Job job)
        {
            try
            {
                var dbJob = await _dbContext.Jobs.SingleOrDefaultAsync(j => j.Id == id);
                if (dbJob == null)
                {
                    return NotFound();
                }

                dbJob.Name = job.Name;
                await _dbContext.SaveChangesAsync();

                return NoContent();
            }
            catch (Exception ex)
            {
                return BadRequest(ex.Message);
            }
        }

        [HttpDelete("{jobId}")]
        public async Task<IActionResult> DeleteStep([FromRoute] Guid jobId)
        {
            try
            {
                _dbContext.Jobs.Remove(await _dbContext.Jobs.SingleAsync(j => j.Id == jobId));
                await _dbContext.SaveChangesAsync();
            }
            catch (Exception ex)
            {
                return Problem(ex.Message, statusCode: StatusCodes.Status500InternalServerError);
            }

            return NoContent();
        }

        [HttpPost("{jobId}/enqueue")]
        public IActionResult Enqueue([FromRoute] Guid jobId, string cron)
        {
            var job = Hangfire.Common.Job.FromExpression(() => Execute(jobId));
            _recurringJobManager.AddOrUpdate(jobId.ToString(), job, cron);

            return NoContent();
        }

        public async Task Execute(Guid jobId)
        {
            var steps = await _dbContext.Steps.Where(s => s.Job.Id == jobId).OrderBy(s => s.Position).ToListAsync();
            foreach (var step in steps)
            {
                switch (step.StepType)
                {
                    case StepType.Ssh:
                    {
                        using var sshClient = new SshClient(step.Server, step.Username, step.Password)
                        {
                            ConnectionInfo = {Timeout = TimeSpan.FromHours(4)}
                        };
                        sshClient.Connect();

                        using var sshCommand = sshClient.CreateCommand(step.Command);
                        sshCommand.CommandTimeout = TimeSpan.FromHours(4);
                        sshCommand.Execute();

                        break;
                    }
                    case StepType.FileTransfer:
                    {
                        using var sftpClient = new SftpClient(step.Server, step.Username, step.Password);
                        sftpClient.Connect();

                        await using var backupFileStream = sftpClient.OpenRead(step.Source);
                        var dirPath = Path.Combine(step.Target, DateTime.Now.ToString("s"));
                        if (!Directory.Exists(dirPath)) Directory.CreateDirectory(dirPath);

                        var filePath = Path.Combine(dirPath, Path.GetFileName(step.Source));
                        await using var fileStream = new FileStream(
                            filePath,
                            FileMode.CreateNew
                        );
                        await backupFileStream.CopyToAsync(fileStream);

                        break;
                    }
                    case StepType.Cleanup:
                    {
                        var files = Directory.GetFiles(step.Target);
                        foreach (var file in files)
                        {
                            var fi = new FileInfo(file);
                            if (fi.CreationTime < DateTime.Now.Subtract(TimeSpan.FromDays(7)))
                            {
                                fi.Delete();
                            }
                        }
                        break;
                    }
                    default:
                        throw new ArgumentOutOfRangeException();
                }
            }
        }
    }
}