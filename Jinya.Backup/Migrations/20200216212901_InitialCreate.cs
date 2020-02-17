using System;
using Microsoft.EntityFrameworkCore.Migrations;

namespace Jinya.Backup.Migrations
{
    public partial class InitialCreate : Migration
    {
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.CreateTable(
                name: "Jobs",
                columns: table => new
                {
                    Id = table.Column<Guid>(nullable: false),
                    Name = table.Column<string>(nullable: false, maxLength: 255)
                },
                constraints: table =>
                {
                    table.UniqueConstraint("UQ_Jobs_Name", x => x.Name);
                    table.PrimaryKey("PK_Jobs", x => x.Id);
                });

            migrationBuilder.CreateTable(
                name: "Steps",
                columns: table => new
                {
                    Id = table.Column<Guid>(nullable: false),
                    Position = table.Column<int>(nullable: false),
                    StepType = table.Column<int>(nullable: false),
                    Command = table.Column<string>(nullable: true),
                    Server = table.Column<string>(nullable: true),
                    Username = table.Column<string>(nullable: true),
                    Password = table.Column<string>(nullable: true),
                    Source = table.Column<string>(nullable: true),
                    Target = table.Column<string>(nullable: true),
                    JobId = table.Column<Guid>(nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_Steps", x => x.Id);
                    table.ForeignKey(
                        name: "FK_Steps_Jobs_JobId",
                        column: x => x.JobId,
                        principalTable: "Jobs",
                        principalColumn: "Id",
                        onDelete: ReferentialAction.Cascade);
                });

            migrationBuilder.CreateIndex(
                name: "IX_Steps_JobId",
                table: "Steps",
                column: "JobId");
        }

        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropTable(
                name: "Steps");

            migrationBuilder.DropTable(
                name: "Jobs");
        }
    }
}