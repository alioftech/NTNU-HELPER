%% 1c: Simulating in rough weather
omega = 0.05; %rad/s
A=45;           %Amplitude
simTime = 3000; % [s]
simout = sim('ship1b','startTime','0','stopTime',sprintf('%d',simTime));
time = simout.get('time');     %get time from workspace
psi = simout.get('psi'); %get psi from workspace

width = 20; % cm
height = 10; % cm
fontsize = 10; %points

set(0,'DefaultTextInterpreter', 'latex') %Interpret (most) text as latex. Since we use set(0,... this is a global setting

%Ploting the figures
figure
plot(time,psi)
xlabel('t [s]')
ylabel('\psi [deg]')
legend('\psi (average heading) ', 'Location','SouthEast')
title('Simulation with $\omega_2=0.05$   (Noise + Waves)')
grid

fig1 = figure(1); %get the figure handle for this specific figure (figure 1), to set figure-specific properties.

fig1.Units = 'centimeters';
fig1.Position = [x y width height];

hgexport(fig1,'1c3w2.eps')


