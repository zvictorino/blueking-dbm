# Generated by Django 3.2.19 on 2024-05-20 03:04

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ("db_meta", "0036_merge_0033_clusterdbhaext_0035_machine_system_info"),
    ]

    operations = [
        migrations.AddField(
            model_name="extraprocessinstance",
            name="bk_instance_id",
            field=models.IntegerField(default=0, help_text="服务实例id，对应cmdb"),
        ),
        migrations.AlterUniqueTogether(
            name="sqlserverdtsinfo",
            unique_together={("ticket_id", "source_cluster_id", "target_cluster_id")},
        ),
    ]
