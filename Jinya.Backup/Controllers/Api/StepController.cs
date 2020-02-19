using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Jinya.Backup.Data;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;

namespace Jinya.Backup.Controllers.Api
{
    [ApiController]
    [Route("api/job/{jobId}/step")]
    public class StepController : Controller
    {
        private readonly ApplicationDbContext _dbContext;

        public StepController(ApplicationDbContext dbContext)
        {
            _dbContext = dbContext;
        }

        [HttpGet]
        public async Task<IEnumerable<Step>> GetSteps([FromRoute] Guid jobId)
        {
            return await _dbContext.Steps.AsNoTracking().Where(s => s.Job.Id == jobId).OrderBy(s => s.Position)
                .ToListAsync();
        }

        [HttpGet("{stepId}")]
        public async Task<Step> GetStep([FromRoute] Guid stepId)
        {
            return await _dbContext.Steps.AsNoTracking().SingleOrDefaultAsync(s => s.Id == stepId);
        }

        [HttpPost]
        public async Task<IActionResult> PostStep([FromRoute] Guid jobId, [FromBody] Step step)
        {
            try
            {
                step.Job = await _dbContext.Jobs.SingleAsync(j => j.Id == jobId);
                await _dbContext.AddAsync(step);
                await _dbContext.SaveChangesAsync();
            }
            catch (Exception ex)
            {
                return BadRequest(ex.Message);
            }

            return NoContent();
        }

        [HttpPut("{stepId}")]
        public async Task<IActionResult> PutStep([FromRoute] Guid stepId, [FromBody] Step step)
        {
            try
            {
                var dbStep = await _dbContext.Steps.SingleOrDefaultAsync(s => s.Id == stepId);
                dbStep.Command = step.Command;
                dbStep.Password = step.Password;
                dbStep.Position = step.Position;
                dbStep.StepType = step.StepType;
                dbStep.Server = step.Server;
                dbStep.Target = step.Target;
                dbStep.Source = step.Source;
                
                await _dbContext.SaveChangesAsync();
            }
            catch (Exception ex)
            {
                return BadRequest(ex.Message);
            }

            return NoContent();
        }

        [HttpPut("{stepId}/move/{position}")]
        public async Task<IActionResult> MoveStep([FromRoute] Guid stepId, [FromRoute] int position)
        {
            try
            {
                var dbStep = await _dbContext.Steps.SingleOrDefaultAsync(s => s.Id == stepId);
                dbStep.Position = position;
                
                await _dbContext.SaveChangesAsync();
            }
            catch (Exception ex)
            {
                return BadRequest(ex.Message);
            }

            return NoContent();
        }

        [HttpDelete("{stepId}")]
        public async Task<IActionResult> DeleteStep([FromRoute] Guid stepId)
        {
            try
            {
                _dbContext.Steps.Remove(await _dbContext.Steps.SingleAsync(s => s.Id == stepId));
                await _dbContext.SaveChangesAsync();
            }
            catch (Exception ex)
            {
                return Problem(detail: ex.Message, statusCode: StatusCodes.Status500InternalServerError);
            }

            return NoContent();
        }
    }
}