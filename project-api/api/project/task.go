package project

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/spxzx/project-api/api/rpc"
	"github.com/spxzx/project-api/pkg/model"
	"github.com/spxzx/project-api/pkg/model/file"
	"github.com/spxzx/project-api/pkg/model/poj"
	"github.com/spxzx/project-api/pkg/model/tasks"
	"github.com/spxzx/project-common/errs"
	"github.com/spxzx/project-common/min"
	"github.com/spxzx/project-common/r"
	"github.com/spxzx/project-common/tms"
	"github.com/spxzx/project-grpc/task/task"
	"net/http"
	"path"
	"strconv"
	"time"
)

type HandlerTask struct{}

func NewTask() *HandlerTask {
	return &HandlerTask{}
}

func (*HandlerTask) getTaskStages(ctx *gin.Context) {
	page := &model.Page{}
	if err := page.Bind(ctx); err != nil {
		r.Fail(ctx, http.StatusBadRequest, "参数格式有误")
		return
	}
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	taskResp, err := rpc.TaskServiceClient.GetTaskStages(c, &task.TaskRpcRequest{
		ProjectCode: ctx.PostForm("projectCode"),
		Page:        page.Page,
		PageSize:    page.PageSize,
	})
	if err != nil {
		e := errs.ParseGrpcError(err)
		r.Fail(ctx, int(e.Code), e.Msg)
		return
	}
	var resp []*tasks.TaskStagesResp
	_ = copier.Copy(&resp, taskResp.List)
	if len(resp) == 0 {
		resp = []*tasks.TaskStagesResp{}
	}
	for _, v := range resp {
		v.TasksLoading = true  //任务加载状态
		v.FixedCreator = false //添加任务按钮定位
		v.ShowTaskCard = false //是否显示创建卡片
		v.Tasks = []int{}
		v.DoneTasks = []int{}
		v.UnDoneTasks = []int{}
	}
	r.Success(ctx, gin.H{
		"list":  resp,
		"total": taskResp.Total,
		"page":  page.Page,
	})
}

func (*HandlerTask) getProjectMember(ctx *gin.Context) {
	page := &model.Page{}
	if err := page.Bind(ctx); err != nil {
		r.Fail(ctx, http.StatusBadRequest, "参数格式有误")
		return
	}
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	pmResp, err := rpc.TaskServiceClient.GetProjectMember(c, &task.TaskRpcRequest{
		ProjectCode: ctx.PostForm("projectCode"),
		Page:        page.Page,
		PageSize:    page.PageSize,
	})
	if err != nil {
		e := errs.ParseGrpcError(err)
		r.Fail(ctx, int(e.Code), e.Msg)
		return
	}
	var resp []*tasks.ProjectMemberResp
	_ = copier.Copy(&resp, pmResp.List)
	if resp == nil {
		resp = []*tasks.ProjectMemberResp{}
	}
	r.Success(ctx, gin.H{
		"list":  resp,
		"total": pmResp.Total,
		"page":  page.Page,
	})
}

func (*HandlerTask) getTaskStageDetailList(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	listResp, err := rpc.TaskServiceClient.GetTaskStageDetailList(c,
		&task.TaskRpcRequest{
			StageCode: ctx.PostForm("stageCode"),
			MemberId:  ctx.GetInt64("memberId"),
		})
	if err != nil {
		e := errs.ParseGrpcError(err)
		r.Fail(ctx, int(e.Code), e.Msg)
		return
	}
	if len(listResp.List) == 0 {
		r.Success(ctx, []*tasks.TaskDisplay{})
		return
	}
	var resp []*tasks.TaskDisplay
	_ = copier.Copy(&resp, listResp.List)
	for _, v := range resp {
		if v.Tags == nil || len(v.Tags) == 0 {
			v.Tags = []int{}
		}
		if v.ChildCount == nil || len(v.ChildCount) == 0 {
			v.ChildCount = []int{}
		}
	}
	r.Success(ctx, resp)
}

func (*HandlerTask) saveTask(ctx *gin.Context) {
	var req *tasks.SaveTaskReq
	if err := ctx.ShouldBind(&req); err != nil {
		r.Fail(ctx, http.StatusBadRequest, "参数格式有误")
		return
	}
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	saveResp, err := rpc.TaskServiceClient.SaveTask(c,
		&task.TaskRpcRequest{
			ProjectCode: req.ProjectCode,
			Name:        req.Name,
			StageCode:   req.StageCode,
			AssignTo:    req.AssignTo,
			MemberId:    ctx.GetInt64("memberId"),
		})
	if err != nil {
		e := errs.ParseGrpcError(err)
		r.Fail(ctx, int(e.Code), e.Msg)
	}
	td := &tasks.TaskDisplay{}
	_ = copier.Copy(td, saveResp)
	if td != nil {
		if td.Tags == nil || len(td.Tags) <= 0 {
			td.Tags = []int{}
		}
		if td.ChildCount == nil || len(td.ChildCount) <= 0 {
			td.ChildCount = []int{}
		}
	}
	r.Success(ctx, td)
}

func (*HandlerTask) moveTask(ctx *gin.Context) {
	var req *tasks.MoveTaskReq
	if err := ctx.ShouldBind(&req); err != nil {
		r.Fail(ctx, http.StatusBadRequest, "参数格式有误")
		return
	}
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if _, err := rpc.TaskServiceClient.MoveTask(c,
		&task.TaskRpcRequest{
			PreTaskCode:  req.PreTaskCode,
			NextTaskCode: req.NextTaskCode,
			ToStageCode:  req.ToStageCode,
		}); err != nil {
		e := errs.ParseGrpcError(err)
		r.Fail(ctx, int(e.Code), e.Msg)
	}
	r.Success(ctx)
}

func (*HandlerTask) getSelfTaskList(ctx *gin.Context) {
	var req *tasks.SelfTaskReq
	if err := ctx.ShouldBind(&req); err != nil {
		r.Fail(ctx, http.StatusBadRequest, "参数格式有误")
		return
	}
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	taskResp, err := rpc.TaskServiceClient.GetSelfTaskList(c, &task.TaskRpcRequest{
		MemberId: ctx.GetInt64("memberId"),
		Page:     int64(req.Page),
		PageSize: int64(req.PageSize),
		TaskType: int32(req.TaskType),
		Type:     int32(req.Type),
	})

	if err != nil {
		e := errs.ParseGrpcError(err)
		r.Fail(ctx, int(e.Code), e.Msg)
		return
	}
	var resp []*tasks.SelfTaskDisplay
	_ = copier.Copy(&resp, taskResp.List)
	for _, v := range resp {
		v.ProjectInfo = tasks.ProjectInfo{
			Name: v.ProjectName,
			Code: v.ProjectCode,
		}
	}
	if len(taskResp.List) == 0 {
		resp = []*tasks.SelfTaskDisplay{}
	}
	r.Success(ctx, gin.H{
		"list":  resp,
		"total": taskResp.Total,
		"page":  req.Page,
	})
}

func (*HandlerTask) readTask(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	taskResp, err := rpc.TaskServiceClient.ReadTask(c, &task.TaskRpcRequest{
		TaskCode: ctx.PostForm("taskCode"),
		MemberId: ctx.GetInt64("memberId"),
	})
	if err != nil {
		e := errs.ParseGrpcError(err)
		r.Fail(ctx, int(e.Code), e.Msg)
		return
	}
	td := &tasks.TaskDisplay{}
	_ = copier.Copy(td, taskResp)
	//_ = copier.Copy(td.Executor, taskResp.Executor)
	//log.Println(taskResp, "\n", td)
	if td != nil {
		if td.Tags == nil || len(td.Tags) <= 0 {
			td.Tags = []int{}
		}
		if td.ChildCount == nil || len(td.Code) <= 0 {
			td.ChildCount = []int{}
		}
	}
	r.Success(ctx, td)
}

func (*HandlerTask) getTaskMember(ctx *gin.Context) {
	page := &model.Page{}
	if err := page.Bind(ctx); err != nil {
		r.Fail(ctx, http.StatusBadRequest, "参数格式有误")
		return
	}
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	tmResp, err := rpc.TaskServiceClient.GetTaskMemberList(c, &task.TaskRpcRequest{
		MemberId: ctx.GetInt64("memberId"),
		Page:     page.Page,
		PageSize: page.PageSize,
		TaskCode: ctx.PostForm("taskCode"),
	})
	if err != nil {
		e := errs.ParseGrpcError(err)
		r.Fail(ctx, int(e.Code), e.Msg)
		return
	}
	var tmList []*tasks.TaskMember
	_ = copier.Copy(&tmList, tmResp.List)
	if tmList == nil || len(tmList) <= 0 {
		tmList = []*tasks.TaskMember{}
	}
	r.Success(ctx, gin.H{
		"list":  tmList,
		"total": tmResp.Total,
		"page":  page.Page,
	})
}

func (*HandlerTask) getTaskLog(ctx *gin.Context) {
	var req *poj.TaskLogReq
	if err := ctx.ShouldBind(&req); err != nil {
		r.Fail(ctx, http.StatusBadRequest, "参数格式有误")
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	taskLogResp, err := rpc.TaskServiceClient.GetTaskLog(c, &task.TaskRpcRequest{
		TaskCode: req.TaskCode,
		MemberId: ctx.GetInt64("memberId"),
		Page:     int64(req.Page),
		PageSize: int64(req.PageSize),
		All:      int32(req.All),
		Comment:  int32(req.Comment),
	})
	if err != nil {
		e := errs.ParseGrpcError(err)
		r.Fail(ctx, int(e.Code), e.Msg)
		return
	}
	var pld []*poj.ProjectLogDisplay
	_ = copier.Copy(&pld, taskLogResp.List)
	if pld == nil || len(pld) <= 0 {
		pld = []*poj.ProjectLogDisplay{}
	}
	r.Success(ctx, gin.H{
		"list":  pld,
		"total": taskLogResp.Total,
		"page":  req.Page,
	})
}

func (*HandlerTask) getTaskWorkTimeList(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	timeResp, err := rpc.TaskServiceClient.GetTaskWorkTimeList(c,
		&task.TaskRpcRequest{
			TaskCode: ctx.PostForm("taskCode"),
			MemberId: ctx.GetInt64("memberId"),
		})
	if err != nil {
		e := errs.ParseGrpcError(err)
		r.Fail(ctx, int(e.Code), e.Msg)
		return
	}
	var tList []*tasks.TaskWorkTime
	_ = copier.Copy(&tList, timeResp.List)
	if tList == nil || len(tList) <= 0 {
		tList = []*tasks.TaskWorkTime{}
	}
	r.Success(ctx, tList)
}

func (*HandlerTask) saveTaskWorkTime(ctx *gin.Context) {
	var req *tasks.SaveTaskWorkTimeReq
	if err := ctx.ShouldBind(&req); err != nil {
		r.Fail(ctx, http.StatusBadRequest, "参数格式有误")
		return
	}
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if _, err := rpc.TaskServiceClient.SaveTaskWorkTimeList(c,
		&task.TaskRpcRequest{
			TaskCode:  req.TaskCode,
			MemberId:  ctx.GetInt64("memberId"),
			Content:   req.Content,
			Num:       int32(req.Num),
			BeginTime: tms.ParseTime(req.BeginTime),
		}); err != nil {
		e := errs.ParseGrpcError(err)
		r.Fail(ctx, int(e.Code), e.Msg)
		return
	}
	r.Success(ctx)
}

func (*HandlerTask) uploadFile(ctx *gin.Context) {
	var req *file.UploadFileReq
	if err := ctx.ShouldBind(&req); err != nil {
		r.Fail(ctx, http.StatusBadRequest, "参数格式有误")
		return
	}
	multipartForm, err := ctx.MultipartForm()
	if err != nil {
		r.Fail(ctx, http.StatusBadRequest, "参数格式有误")
		return
	}
	f := multipartForm.File
	key := "pmproject/" + req.Filename
	minioClient, err := min.New(
		"localhost:9009",
		"VGsM70md0bxAX5DE",
		"CCFGeMeWHwrpHt2bJ1UW8HzaA0aLzvjE",
		false,
	)
	if err != nil {
		r.Fail(ctx, 999, "连接MinIO失败")
		return
	}
	if req.TotalChunks == 1 { // 代表不分片，直接上传
		/*p := "upload/" + req.ProjectCode + "/" + req.TaskCode + "/" + tms.FormatYMD(time.Now())
		if !fs.IsExist(p) {
			_ = os.MkdirAll(p, os.ModePerm)
		}
		dst := p + "/" + req.Filename
		key = dst
		header := f["file"][0]
		if err := ctx.SaveUploadedFile(header, dst); err != nil {
			r.Fail(ctx, 9999, err.Error())
			return
		}*/
		open, _ := f["file"][0].Open()
		defer open.Close()
		buf := make([]byte, req.CurrentChunkSize)
		_, _ = open.Read(buf)

		_, err = minioClient.Put(
			context.Background(),
			"pmproject",
			req.Filename,
			buf,
			int64(req.TotalSize),
			f["file"][0].Header.Get("Content-Type"),
		)
		if err != nil {
			r.Fail(ctx, 9999, err.Error())
			return
		}
	}
	if req.TotalChunks > 1 {
		/*p := "upload/" + req.ProjectCode + "/" + req.TaskCode + "/" + tms.FormatYMD(time.Now())
		if !fs.IsExist(p) {
			_ = os.MkdirAll(p, os.ModePerm)
		}
		fileName := p + "/" + req.Identifier
		openFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
		if err != nil {
			r.Fail(ctx, 9999, err.Error())
			return
		}*/
		open, err := f["file"][0].Open()
		defer open.Close()
		buf := make([]byte, req.CurrentChunkSize)
		_, _ = open.Read(buf)
		/*_, _ = openFile.Write(buf)
		_ = openFile.Close()
		newpath := p + "/" + req.Filename
		key = newpath*/
		fInt := strconv.FormatInt(int64(req.ChunkNumber), 10)

		if _, err = minioClient.Put(
			context.Background(),
			"pmproject",
			req.Filename+"_"+fInt,
			buf,
			int64(req.CurrentChunkSize),
			f["file"][0].Header.Get("Content-Type"),
		); err != nil {
			r.Fail(ctx, 9999, err.Error())
			return
		}
		if req.TotalChunks == req.ChunkNumber {
			//最后一块 重命名文件名
			//_ = os.Rename(fileName, newpath)
			// 合并文件
			if _, err := minioClient.Compose(
				context.Background(),
				"pmproject",
				req.Filename,
				req.TotalChunks,
			); err != nil {
				r.Fail(ctx, 9999, err.Error())
				return
			}
		}
	}
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	fileUrl := "http://localhost:9009/" + key
	if req.TotalChunks == req.ChunkNumber {
		_, err = rpc.TaskServiceClient.UploadFile(c, &task.TaskFileRequest{
			PathName:         key,
			FileName:         req.Filename,
			Extension:        path.Ext(key),
			Size:             int64(req.TotalSize),
			ProjectCode:      req.ProjectCode,
			TaskCode:         req.TaskCode,
			OrganizationCode: ctx.GetString("organizationCode"),
			FileUrl:          fileUrl,
			FileType:         f["file"][0].Header.Get("Content-Type"),
			MemberId:         ctx.GetInt64("memberId"),
		})
		if err != nil {
			e := errs.ParseGrpcError(err)
			r.Fail(ctx, int(e.Code), e.Msg)
			return
		}
	}
	r.Success(ctx, gin.H{
		"file":        key,
		"hash":        "",
		"key":         key,
		"url":         "http://localhost:9009/" + key,
		"projectName": req.ProjectName,
	})
}

func (*HandlerTask) getTaskSourceLink(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	sourceResp, err := rpc.TaskServiceClient.GetTaskSources(c, &task.TaskRpcRequest{TaskCode: ctx.PostForm("taskCode")})
	if err != nil {
		e := errs.ParseGrpcError(err)
		r.Fail(ctx, int(e.Code), e.Msg)
		return
	}
	var slList []*file.SourceLink
	_ = copier.Copy(&slList, sourceResp.List)
	if slList == nil || len(slList) <= 0 {
		slList = []*file.SourceLink{}
	}
	r.Success(ctx, slList)
}

func (*HandlerTask) createComment(ctx *gin.Context) {
	var req *poj.CommentReq
	if err := ctx.ShouldBind(&req); err != nil {
		r.Fail(ctx, http.StatusBadRequest, "参数格式有误")
		return
	}
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if _, err := rpc.TaskServiceClient.CreateComment(c, &task.TaskRpcRequest{
		TaskCode:       req.TaskCode,
		CommentContent: req.Comment,
		Mentions:       req.Mentions,
		MemberId:       ctx.GetInt64("memberId"),
	}); err != nil {
		e := errs.ParseGrpcError(err)
		r.Fail(ctx, int(e.Code), e.Msg)
		return
	}
	r.Success(ctx)
}

//if err != nil {
//e := errs.ParseGrpcError(err)
//r.Fail(ctx, int(e.Code), e.Msg)
//return
//}
