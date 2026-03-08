package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &TribalResourceResource{}
var _ resource.ResourceWithImportState = &TribalResourceResource{}

type TribalResourceResource struct {
	client *TribalClient
}

type TribalResourceModel struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	DRI                    types.String `tfsdk:"dri"`
	Type                   types.String `tfsdk:"type"`
	ExpirationDate         types.String `tfsdk:"expiration_date"`
	Purpose                types.String `tfsdk:"purpose"`
	GenerationInstructions types.String `tfsdk:"generation_instructions"`
	SecretManagerLink      types.String `tfsdk:"secret_manager_link"`
	SlackWebhook           types.String `tfsdk:"slack_webhook"`
	PublicKeyPEM           types.String `tfsdk:"public_key_pem"`
	CreatedAt              types.String `tfsdk:"created_at"`
	UpdatedAt              types.String `tfsdk:"updated_at"`
}

func NewTribalResourceResource() resource.Resource {
	return &TribalResourceResource{}
}

func (r *TribalResourceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource"
}

func (r *TribalResourceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Tribal tracked resource (certificate, API key, SSH key, etc.).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Numeric identifier of the resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the resource.",
			},
			"dri": schema.StringAttribute{
				Required:    true,
				Description: "Directly Responsible Individual for this resource.",
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "Type of resource: Certificate, API Key, SSH Key, or Other.",
			},
			"expiration_date": schema.StringAttribute{
				Required:    true,
				Description: "Expiration date in YYYY-MM-DD format.",
			},
			"purpose": schema.StringAttribute{
				Required:    true,
				Description: "Purpose/description of the resource.",
			},
			"generation_instructions": schema.StringAttribute{
				Required:    true,
				Description: "Instructions for generating/renewing this resource.",
			},
			"secret_manager_link": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "URL link to the secret in a secret manager.",
			},
			"slack_webhook": schema.StringAttribute{
				Required:    true,
				Description: "Slack webhook URL for expiration notifications.",
			},
			"public_key_pem": schema.StringAttribute{
				Computed:    true,
				Description: "PEM-encoded public certificate (if uploaded).",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the resource was created.",
			},
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the resource was last updated.",
			},
		},
	}
}

func (r *TribalResourceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*TribalClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *TribalClient, got: %T", req.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *TribalResourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TribalResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := ResourceRequest{
		Name:                   plan.Name.ValueString(),
		DRI:                    plan.DRI.ValueString(),
		Type:                   plan.Type.ValueString(),
		ExpirationDate:         plan.ExpirationDate.ValueString(),
		Purpose:                plan.Purpose.ValueString(),
		GenerationInstructions: plan.GenerationInstructions.ValueString(),
		SecretManagerLink:      plan.SecretManagerLink.ValueString(),
		SlackWebhook:           plan.SlackWebhook.ValueString(),
	}

	apiResp, err := r.client.CreateResource(createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Resource", err.Error())
		return
	}

	plan.ID = types.StringValue(strconv.Itoa(apiResp.ID))
	plan.Name = types.StringValue(apiResp.Name)
	plan.DRI = types.StringValue(apiResp.DRI)
	plan.Type = types.StringValue(apiResp.Type)
	plan.ExpirationDate = types.StringValue(apiResp.ExpirationDate)
	plan.Purpose = types.StringValue(apiResp.Purpose)
	plan.GenerationInstructions = types.StringValue(apiResp.GenerationInstructions)
	plan.SecretManagerLink = types.StringValue(apiResp.SecretManagerLink)
	plan.SlackWebhook = types.StringValue(apiResp.SlackWebhook)
	plan.PublicKeyPEM = types.StringValue(apiResp.PublicKeyPEM)
	plan.CreatedAt = types.StringValue(apiResp.CreatedAt)
	plan.UpdatedAt = types.StringValue(apiResp.UpdatedAt)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *TribalResourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TribalResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid Resource ID", err.Error())
		return
	}

	apiResp, err := r.client.GetResource(id)
	if err != nil {
		if strings.Contains(err.Error(), "API error 404") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error Reading Resource", err.Error())
		return
	}

	state.ID = types.StringValue(strconv.Itoa(apiResp.ID))
	state.Name = types.StringValue(apiResp.Name)
	state.DRI = types.StringValue(apiResp.DRI)
	state.Type = types.StringValue(apiResp.Type)
	state.ExpirationDate = types.StringValue(apiResp.ExpirationDate)
	state.Purpose = types.StringValue(apiResp.Purpose)
	state.GenerationInstructions = types.StringValue(apiResp.GenerationInstructions)
	state.SecretManagerLink = types.StringValue(apiResp.SecretManagerLink)
	state.SlackWebhook = types.StringValue(apiResp.SlackWebhook)
	state.PublicKeyPEM = types.StringValue(apiResp.PublicKeyPEM)
	state.CreatedAt = types.StringValue(apiResp.CreatedAt)
	state.UpdatedAt = types.StringValue(apiResp.UpdatedAt)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *TribalResourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TribalResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	var state TribalResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid Resource ID", err.Error())
		return
	}

	updateReq := ResourceRequest{
		Name:                   plan.Name.ValueString(),
		DRI:                    plan.DRI.ValueString(),
		Type:                   plan.Type.ValueString(),
		ExpirationDate:         plan.ExpirationDate.ValueString(),
		Purpose:                plan.Purpose.ValueString(),
		GenerationInstructions: plan.GenerationInstructions.ValueString(),
		SecretManagerLink:      plan.SecretManagerLink.ValueString(),
		SlackWebhook:           plan.SlackWebhook.ValueString(),
	}

	apiResp, err := r.client.UpdateResource(id, updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Resource", err.Error())
		return
	}

	plan.ID = types.StringValue(strconv.Itoa(apiResp.ID))
	plan.Name = types.StringValue(apiResp.Name)
	plan.DRI = types.StringValue(apiResp.DRI)
	plan.Type = types.StringValue(apiResp.Type)
	plan.ExpirationDate = types.StringValue(apiResp.ExpirationDate)
	plan.Purpose = types.StringValue(apiResp.Purpose)
	plan.GenerationInstructions = types.StringValue(apiResp.GenerationInstructions)
	plan.SecretManagerLink = types.StringValue(apiResp.SecretManagerLink)
	plan.SlackWebhook = types.StringValue(apiResp.SlackWebhook)
	plan.PublicKeyPEM = types.StringValue(apiResp.PublicKeyPEM)
	plan.CreatedAt = types.StringValue(apiResp.CreatedAt)
	plan.UpdatedAt = types.StringValue(apiResp.UpdatedAt)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *TribalResourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TribalResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid Resource ID", err.Error())
		return
	}

	if err := r.client.DeleteResource(id); err != nil {
		if strings.Contains(err.Error(), "API error 404") {
			return
		}
		resp.Diagnostics.AddError("Error Deleting Resource", err.Error())
	}
}

func (r *TribalResourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID
	if _, err := strconv.Atoi(id); err != nil {
		resp.Diagnostics.AddError("Invalid Import ID", "Import ID must be a numeric resource ID")
		return
	}

	apiResp, err := r.client.GetResource(mustAtoi(id))
	if err != nil {
		resp.Diagnostics.AddError("Error Importing Resource", err.Error())
		return
	}

	state := TribalResourceModel{
		ID:                     types.StringValue(strconv.Itoa(apiResp.ID)),
		Name:                   types.StringValue(apiResp.Name),
		DRI:                    types.StringValue(apiResp.DRI),
		Type:                   types.StringValue(apiResp.Type),
		ExpirationDate:         types.StringValue(apiResp.ExpirationDate),
		Purpose:                types.StringValue(apiResp.Purpose),
		GenerationInstructions: types.StringValue(apiResp.GenerationInstructions),
		SecretManagerLink:      types.StringValue(apiResp.SecretManagerLink),
		SlackWebhook:           types.StringValue(apiResp.SlackWebhook),
		PublicKeyPEM:           types.StringValue(apiResp.PublicKeyPEM),
		CreatedAt:              types.StringValue(apiResp.CreatedAt),
		UpdatedAt:              types.StringValue(apiResp.UpdatedAt),
	}

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func mustAtoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
